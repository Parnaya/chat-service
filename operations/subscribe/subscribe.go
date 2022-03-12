package subscribe

import (
	"chat.service/operations/entity"
	"encoding/json"
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
}

var (
	server = Server{
		Clients: make(map[*websocket.Conn]*Client),
	}
)

func handleWebSocketMessage(
	connection *websocket.Conn,
	handleCreate func(entity entity.Entity),
	sendUpdateChan chan *entity.Entity,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			return
		}

		item := new(entity.Entity)
		json.Unmarshal(message, item)
		handleCreate(*item)

		sendUpdateChan <- item
	}
}

func publishMessageToWebSocket(
	connection *websocket.Conn,
	receiveUpdateChan chan *entity.Entity,
) {
	for it := range receiveUpdateChan {
		rawMessage, _ := json.Marshal(it)
		connection.WriteMessage(websocket.TextMessage, rawMessage)
	}
}

func Echo(handleCreate func(entity entity.Entity)) echo.HandlerFunc {
	return func(config echo.Context) error {
		connection, err := upgrader.Upgrade(config.Response(), config.Request(), nil)

		if err != nil {
			return err
		}

		receiveChannel := make(chan *entity.Entity)

		client := &Client{receiveChannel}
		server.Clients[connection] = client

		var wg sync.WaitGroup
		wg.Add(1)
		go handleWebSocketMessage(connection, handleCreate, receiveChannel, &wg)
		go publishMessageToWebSocket(connection, receiveChannel)

		defer close(receiveChannel)
		defer delete(server.Clients, connection)
		defer connection.Close()

		wg.Wait()

		return nil
	}
}
