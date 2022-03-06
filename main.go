package main

import (
	"chat.service/configuration"
	"chat.service/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	configuration.ShouldParseViperConfig()
	couchbaseConfig := configuration.ShouldParseCouchbaseConfig()
	_ = database.ShouldGetCluster(couchbaseConfig)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
