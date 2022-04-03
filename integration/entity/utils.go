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

func reverse(items []interface{}) []interface{} {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return items
}

func ter(is bool, a interface{}, b interface{}) interface{} {
	if is {
		return a
	}

	return b
}
