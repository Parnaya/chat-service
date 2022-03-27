package mapper

import (
	"chat.service/model"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

import _ "github.com/santhosh-tekuri/jsonschema/v5/httploader"

type JsonSocketRequest struct {
	Id        string                     `json:"id"`
	CreatedAt string                     `json:"createdAt"`
	Messages  []JsonSocketRequestMessage `json:"messages"`
}

type JsonSocketRequestMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type CreateEntityMessage struct {
	Id   string      `json:"id"`
	Tags []string    `json:"tags"`
	Data interface{} `json:"data"`
}

func JsonSocketRequestMapper(schema *jsonschema.Schema) func(messageBytes []byte) *model.SocketRequest {
	return func(messageBytes []byte) *model.SocketRequest {
		var jsonObject map[string]interface{}

		if err := json.Unmarshal(messageBytes, &jsonObject); err != nil {
			fmt.Println(err)
			return nil
		}

		if err := schema.Validate(jsonObject); err != nil {
			fmt.Println(err)
			return nil
		}

		jsonRequest := JsonSocketRequest{}
		if err := mapstructure.Decode(jsonObject, &jsonRequest); err != nil {
			return nil
		}

		request := new(model.SocketRequest)
		request.Messages = make([]model.SocketRequestMessage, len(jsonRequest.Messages))

		for messageIndex, message := range jsonRequest.Messages {
			switch message.Type {
			case "insert":
				entityCreate := CreateEntityMessage{}
				if err := mapstructure.Decode(message.Data, &entityCreate); err != nil {
					break
				}

				item := new(model.Entity)
				if err := item.Id.UnmarshalText([]byte(entityCreate.Id)); err != nil {
					break
				}

				item.Tags = entityCreate.Tags

				data, ok := entityCreate.Data.(map[string]interface{})

				if !ok {
					return nil
				}

				item.Data = data

				request.Messages[messageIndex] = model.SocketRequestMessage{
					RequestType: model.Create,
					Data:        item,
				}

				break
			case "filters":
				request.Messages[messageIndex] = model.SocketRequestMessage{
					RequestType: model.Filters,
					Data:        message.Data,
				}

				break
			}

		}

		return request
	}
}
