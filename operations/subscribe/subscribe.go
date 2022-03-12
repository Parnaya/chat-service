package subscribe

import (
	"chat.service/operations/entity"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
	receiveChannel chan *entity.Entity
	filters        []string
}

type Request map[string]interface{} // { 'filters': ['', '', ...], 'entity_create': {  } }

var (
	server = Server{
		Clients: make(map[*websocket.Conn]*Client),
	}
)

func handleWebSocketMessage(
	connection *websocket.Conn,
	handleCreate func(entity entity.Entity),
	client Client,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			return
		}

		request := new(Request)
		json.Unmarshal(message, request)

		for key, element := range *request {
			switch key {
			case "filters":
				hh, _ := json.Marshal(element)
				filters := new([]string)
				json.Unmarshal(hh, filters)
				client.filters = *filters
			case "entity_create":
				hh, _ := json.Marshal(element)
				newEntity := new(entity.Entity)
				json.Unmarshal(hh, newEntity)
				handleCreate(*newEntity)
				client.receiveChannel <- newEntity
			case "entity_update":
				fmt.Print("not handle")
			case "entity_delete":
				fmt.Print("not handle")
			}
		}
	}
}

func publishMessageToWebSocket(
	connection *websocket.Conn,
	client Client,
) {
	for it := range client.receiveChannel {
		isMatch := containsAll(client.filters, it.Tags)
		if isMatch && len(client.filters) > 0 {
			rawMessage, _ := json.Marshal(it)
			connection.WriteMessage(websocket.TextMessage, rawMessage)
		}
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

func Echo(handleCreate func(entity entity.Entity)) echo.HandlerFunc {
	return func(config echo.Context) error {
		connection, err := upgrader.Upgrade(config.Response(), config.Request(), nil)

		if err != nil {
			return err
		}

		receiveChannel := make(chan *entity.Entity)

		client := &Client{receiveChannel, []string{}}
		server.Clients[connection] = client

		var wg sync.WaitGroup
		wg.Add(1)
		go handleWebSocketMessage(connection, handleCreate, *client, &wg)
		go publishMessageToWebSocket(connection, *client)

		defer close(receiveChannel)
		defer delete(server.Clients, connection)
		defer connection.Close()

		wg.Wait()

		return nil
	}
}
