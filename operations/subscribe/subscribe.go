package subscribe

import (
	"chat.service/integration/entity"
	"chat.service/model"
	"chat.service/operations/log"
	"encoding/json"
	"fmt"
	"github.com/5anthosh/chili"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	sql            string
}

var (
	server = Server{
		Clients: make(map[*websocket.Conn]*Client),
	}
)

type SubscribeOperationSettings struct {
	SocketRequestMapper func(message []byte) *model.SocketRequest
	Entity              entity.Entity
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

				settings.Entity.Create(entity)

				for _, serverClient := range server.Clients {
					if !serverClient.isMatch(entity.Tags) {
						continue
					}

					serverClient.receiveChannel <- messageBytes
				}
				break
			case model.Filters:
				expr, ok := it.Data.(string)
				if !ok {
					break
				}

				// TODO: сделать нормальную регулярку
				def := regexp.MustCompile(`[ |&|\||!]+`).Split(expr, -1)

				sql := expr
				sql = strings.Replace(sql, "&&", "AND", -1)
				sql = strings.Replace(sql, "||", "OR", -1)

				for _, tag := range def {
					strings.Replace(tag, tag, "tag = `"+tag+"`", -1)
				}

				fmt.Println("AND ANY tag IN `tags` SATISFIES " + sql)

				//server.Clients[connection].sql =

				server.Clients[connection].isMatch = func(tags []string) bool {
					next := expr

					for i, tags := range [][]string{tags, def} {
						for _, tag := range tags {
							next = strings.Replace(next, tag, strconv.FormatBool(i == 0), -1)
						}
					}

					result := log.Proxy(chili.Eval(next, map[string]interface{}{}))

					if result == nil {
						return false
					}

					return result.(bool)
				}
				break

			case model.Fetch:
				params, ok := it.Data.(*entity.GetParams)

				params.Filters = server.Clients[connection].sql

				if !ok {
					break
				}

				items := settings.Entity.Get(params)

				if len(items) == 0 {
					break
				}

				var messages []interface{}

				for _, v := range settings.Entity.Get(params) {
					// TODO: validate
					item := make(map[string]interface{})

					item["type"] = "insert"
					item["data"] = v

					messages = append(messages, item)
				}

				next := make(map[string]interface{})
				next["id"] = ""
				next["messages"] = messages

				nextBytes := log.Proxy(json.Marshal(next)).([]byte)

				server.Clients[connection].receiveChannel <- nextBytes

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
		connection.WriteMessage(websocket.TextMessage, it)
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

		client := &Client{receiveChannel, isMatchDef, ""}

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
