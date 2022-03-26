package subscribe

import (
	"chat.service/model"
	"github.com/5anthosh/chili"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
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
	receiveChannel chan []byte
	isMatch        func([]string) bool
}

var (
	server = Server{
		Clients: make(map[*websocket.Conn]*Client),
	}
)

type SubscribeOperationSettings struct {
	SocketRequestMapper func(message []byte) *model.SocketRequest
	HandleEntityCreate  func(entity *model.Entity)
}

func handleWebSocketMessage(
	settings *SubscribeOperationSettings,
	connection *websocket.Conn,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		mt, messageBytes, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			return
		}

		request := settings.SocketRequestMapper(messageBytes)
		if request == nil {
			continue
		}

		for _, it := range request.Messages {
			switch it.RequestType {
			case model.Create:
				entity, ok := it.Data.(*model.Entity)
				if !ok {
					break
				}

				settings.HandleEntityCreate(entity)

				for _, serverClient := range server.Clients {
					if !serverClient.isMatch(entity.Tags) {
						continue
					}

					serverClient.receiveChannel <- messageBytes
				}
				break
			case model.Filters:
				expr := it.Data.(string)

				def := regexp.MustCompile(`-|\+|&&|\|\||\(|\)`).Split(expr, -1)

				server.Clients[connection].isMatch = func(tags []string) bool {
					values := map[string]interface{}{}
					for i, tags := range [][]string{def, tags} {
						for _, tag := range tags {
							values[tag] = i
						}
					}

					result, err := chili.Eval(expr, values)

					if err != nil {
						panic(err)
					}

					return result.(bool)
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

		connection.WriteMessage(websocket.TextMessage, it)
		//}
	}
}

func isMatchDef([]string) bool {
	return false
}

func OpenWebSocketConnection(operationSettings *SubscribeOperationSettings) echo.HandlerFunc {
	return func(config echo.Context) error {
		connection, err := upgrader.Upgrade(config.Response(), config.Request(), nil)

		if err != nil {
			return err
		}

		receiveChannel := make(chan []byte)

		client := &Client{receiveChannel, isMatchDef}

		server.Clients[connection] = client

		var wg sync.WaitGroup
		wg.Add(1)
		go handleWebSocketMessage(operationSettings, connection, &wg)
		go publishMessageToWebSocket(connection, client)

		defer close(receiveChannel)
		defer delete(server.Clients, connection)
		defer connection.Close()

		wg.Wait()

		return nil
	}
}
