package rest

import (
	"github.com/labstack/echo/v4"
	v1 "github.com/neyrzx/youmusic/internal/delivery/rest/v1"
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
// @BasePath /api/v1
func InitAPI(e *echo.Echo, ts v1.TracksService) {
	api := e.Group("api/v1")

	tracksGroup := api.Group("/tracks")
	v1.NewTracksHandlers(tracksGroup, ts)
}
