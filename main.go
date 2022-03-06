package main

import (
	"chat.service/configuration"
	"chat.service/database"
	"chat.service/operations/entity"
	"chat.service/operations/subscribe"
	"github.com/labstack/echo/v4"
)

func main() {
	configuration.ShouldParseViperConfig()
	couchbaseConfig := configuration.ShouldParseCouchbaseConfig()
	_ = database.ShouldGetCluster(couchbaseConfig)

	e := echo.New()

	e.Static("/", "./public")

	e.GET("/ws", subscribe.Echo(entity.Create))

	e.Logger.Fatal(e.Start(":1323"))
}
