package subscribe

import (
	"chat.service/operations/entity"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type Server struct {
	Clients map[*websocket.Conn]bool
}

var (
	server = Server{
		make(map[*websocket.Conn]bool),
	}
)

func Echo(handleCreate func(entity entity.Entity)) echo.HandlerFunc {
	return func(config echo.Context) error {
		connection, _ := upgrader.Upgrade(config.Response(), config.Request(), nil)

		// TODO: отправлять в 1 соединении
		for _, v := range entity.Database {
			item, _ := json.Marshal(v)

			connection.WriteMessage(websocket.TextMessage, item)
		}

		server.Clients[connection] = true
		defer delete(server.Clients, connection)
		defer connection.Close()

		for {
			mt, message, err := connection.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				return nil
			}

			item := new(entity.Entity)

			json.Unmarshal(message, item)

			handleCreate(*item)

			for conn := range server.Clients {
				// TODO: filtration
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}
}
