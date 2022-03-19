package entity

import (
	"chat.service/model"
	"fmt"
	"github.com/couchbase/gocb/v2"
)

func CouchbaseCreate(collection *gocb.Collection) func(entity *model.Entity) {
	return func(entity *model.Entity) {
		if _, err := collection.Insert(entity.Id.String(), entity, nil); err != nil {
			fmt.Errorf("[Couchbase] Ошибка во время вставки entity: %s", err)
		}
	}
}

func CouchbaseUpdate(collection *gocb.Collection) func(entity *model.Entity) {
	return func(entity *model.Entity) {

	}
}

func CouchbaseDelete(collection *gocb.Collection) func(entity *model.Entity) {
	return func(entity *model.Entity) {

	}
}
