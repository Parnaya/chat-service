package subscribe

import (
	"fmt"
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

func Echo(config echo.Context) error {
	w := config.Response()
	r := config.Request()

	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	server.Clients[connection] = true
	defer delete(server.Clients, connection)

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			return nil
		}

		fmt.Println(string(message))

		// TODO: operation message create
	}
}

func Write(message []byte, filters []string) {
	for conn := range server.Clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
