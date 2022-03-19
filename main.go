package main

import (
	"chat.service/configuration"
	"chat.service/database"
	"chat.service/integration/entity"
	"chat.service/operations/subscribe"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func main() {
	configuration.ShouldParseViperConfig()
	couchbaseConfig := configuration.ShouldParseCouchbaseConfig()
	cluster := database.ShouldGetCluster(couchbaseConfig)
	if err := cluster.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}
	bucket := cluster.Bucket("woop")
	if err := bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}
	entityCollection := bucket.DefaultCollection()

	e := echo.New()

	e.Use(middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:   "public",
			Index:  "index.html",
			Browse: false,
			HTML5:  true,
		},
	))

	e.GET("/ws", subscribe.OpenWebSocketConnection(entity.CouchbaseCreate(entityCollection)))

	e.Logger.Fatal(e.Start(":1323"))
}
