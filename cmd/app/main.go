package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/neyrzx/youmusic/cmd/app/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title YouMusic
// @version 0.0.1
// @description Это проект был разработан в рамках тестового задания от EffectiveMobile
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email khorev.valeriy@yandex.ru

// @license.name MIT
// @license.url https://github.com/Neyrzx/YouMusic?tab=MIT-1-ov-file

// @host localhost:9090
// @BasePath /v1
func main() {
	e := echo.New()
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":9090"))
}
