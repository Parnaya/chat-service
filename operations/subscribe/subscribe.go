package subscribe

import (
	woop "chat.service/gen/github.com/Parnaya/woop-common"
	"chat.service/model"
	"chat.service/utils"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type Server struct {
	Clients map[*websocket.Conn]*Client
}

type Client struct {
	receiveChannel chan *woop.WoopSocketMessage
	filters        []string
}

var (
	server = Server{
		Clients: make(map[*websocket.Conn]*Client),
	}
)

func handleWebSocketMessage(
	connection *websocket.Conn,
	handleCreate func(entity *model.Entity),
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		mt, messageBytes, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			return
		}

		protoMessage := &woop.WoopSocketMessage{}
		if err := proto.Unmarshal(messageBytes, protoMessage); err != nil {
			return
		}

		for _, wrapper := range protoMessage.GetWrapper() {
			switch message := wrapper.GetMessage().(type) {
			case *woop.MessageWrapper_EntityCreate:
				entityCreate := message.EntityCreate
				item := new(model.Entity)
				if err := item.Id.UnmarshalBinary(entityCreate.Id); err != nil {
					break
				}

				item.Tags = make([]model.Tag, len(entityCreate.GetTags()))
				for i, protoTag := range entityCreate.GetTags() {
					tag := model.Tag{}
					if err := item.Id.UnmarshalBinary(protoTag.Id); err != nil {
						break
					}

					if tag.Data, err = utils.ProtobufStructToMap(protoTag.Data); err != nil {
						break
					}

					item.Tags[i] = tag
				}

				handleCreate(item)

				for _, serverClient := range server.Clients {
					serverClient.receiveChannel <- protoMessage
				}
				break
			}
		}

	}
}

func publishMessageToWebSocket(
	connection *websocket.Conn,
	client *Client,
) {
	for it := range client.receiveChannel {
		//isMatch := containsAll(it.Tags, client.filters)
		//if isMatch && len(client.filters) > 0 {

		//protoSocketMessage := new(woop.WoopSocketMessage)
		//
		//messageId, _ := uuid.NewUUID()
		//messageIdBytes, _ := messageId.MarshalBinary()
		//
		//protoSocketMessage.Id = messageIdBytes
		//protoSocketMessage.CreatedAt = timestamppb.Now()

		rawMessage, _ := json.Marshal(it)
		connection.WriteMessage(websocket.BinaryMessage, rawMessage)
		//}
	}
}

func containsAll(container []string, elements []string) bool {
	is := true
	for _, elem := range elements {
		is = is && contains(container, elem)
	}
	return is
}

func contains(container []string, element string) bool {
	for _, entry := range container {
		if entry == element {
			return true
		}
	}
	return false
}

func OpenWebSocketConnection(handleCreate func(entity *model.Entity)) echo.HandlerFunc {
	return func(config echo.Context) error {
		connection, err := upgrader.Upgrade(config.Response(), config.Request(), nil)

		if err != nil {
			return err
		}

		receiveChannel := make(chan *woop.WoopSocketMessage)

		client := &Client{receiveChannel, make([]string, 0)}
		server.Clients[connection] = client

		var wg sync.WaitGroup
		wg.Add(1)
		go handleWebSocketMessage(connection, handleCreate, &wg)
		go publishMessageToWebSocket(connection, client)

		defer close(receiveChannel)
		defer delete(server.Clients, connection)
		defer connection.Close()

		wg.Wait()

		return nil
	}
}
