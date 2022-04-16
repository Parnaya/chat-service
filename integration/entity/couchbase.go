package entity

import (
	"chat.service/model"
	"chat.service/operations/log"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strings"
	"time"
)

type GetParams struct {
	After   interface{} `json:"after"`
	Size    interface{} `json:"size"`
	Filters interface{} `json:"filters"`
}

type Entity struct {
	Create func(entity *model.Entity)
	Update func(entity *model.Entity)
	Delete func(entity *model.Entity)
	Get    func(params *GetParams) []interface{}
}

func couchbaseGet(cluster *gocb.Cluster) func(params *GetParams) []interface{} {
	return func(params *GetParams) []interface{} {

		args := make(map[string]interface{})

		args["after"] = ter(params.After == nil, "", params.After)
		args["size"] = ter(params.Size == nil, 20, params.Size)

		sql := "SELECT * FROM `woop` WHERE `createdAt` < $after $filters ORDER BY `createdAt` DESC LIMIT $size"

		// TODO: разрабы коучбейс долбаебы, реплейс не могут сделать
		sql = strings.Replace(sql, "$filters", params.Filters.(string), -1)

		rows := log.Proxy(
			cluster.Query(sql, &gocb.QueryOptions{NamedParameters: args}),
		).(*gocb.QueryResult)

		var item map[string]interface{}
		var items []interface{}

		for rows.Next() {
			rows.Row(&item)

			items = append(items, item["woop"])
		}

		return reverse(items)
	}
}

func couchbaseCreate(collection *gocb.Collection) func(entity *model.Entity) {
	return func(entity *model.Entity) {
		if _, err := collection.Insert(entity.Id.String(), entity, nil); err != nil {
			fmt.Errorf("[Couchbase] Ошибка во время вставки entity: %s", err)
		}
	}
}

func couchbaseUpdate(collection *gocb.Collection) func(entity *model.Entity) {
	return func(entity *model.Entity) {

	}
}

func couchbaseDelete(collection *gocb.Collection) func(entity *model.Entity) {
	return func(entity *model.Entity) {

	}
}

func Handlers(cluster *gocb.Cluster) Entity {
	name := "woop"

	bucket := cluster.Bucket(name)
	indexes := cluster.QueryIndexes()

	if err := bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	if err := indexes.DropPrimaryIndex(name, nil); err != nil {
		fmt.Println("err", err)
	}

	if err := indexes.CreatePrimaryIndex(name, nil); err != nil {
		fmt.Println("err", err)
	}

	collection := bucket.DefaultCollection()

	return Entity{
		Get:    couchbaseGet(cluster),
		Create: couchbaseCreate(collection),
		Update: couchbaseUpdate(collection),
		Delete: couchbaseDelete(collection),
	}
}
