package entity

import (
	"chat.service/model"
	"github.com/google/uuid"
)

func with(id uuid.UUID, block func(i int, entity *model.Entity)) {
	for i, v := range Database {
		if v.Id == id {
			block(i, v)
		}
	}
}
