package mapper

import (
	woop "chat.service/gen/github.com/Parnaya/woop-common"
	"chat.service/model"
	"google.golang.org/protobuf/proto"
)

func ProtobufSocketRequestMapper(messageBytes []byte) *model.SocketRequest {
	protoMessage := &woop.WoopSocketMessage{}
	if err := proto.Unmarshal(messageBytes, protoMessage); err != nil {
		return nil
	}

	request := new(model.SocketRequest)
	request.Messages = make([]model.SocketRequestMessage, len(protoMessage.GetWrapper()))

	for wrapperIndex, wrapper := range protoMessage.GetWrapper() {
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

				tag.Data = protoTag.Data.AsMap()

				item.Tags[i] = tag
			}

			item.Data = entityCreate.Data.AsMap()

			request.Messages[wrapperIndex] = model.SocketRequestMessage{
				RequestType: model.Create,
				Data:        item,
			}

			break
		}
	}

	return request
}
