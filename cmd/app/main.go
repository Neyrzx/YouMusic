package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/neyrzx/youmusic/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":9090"))
}
