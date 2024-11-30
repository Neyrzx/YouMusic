package v1

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/neyrzx/youmusic/internal/domain/entities"
)

type TracksService interface {
	Create(ctx context.Context, track entities.TrackCreate) error
	GetByID(ctx context.Context, ID int) (entities.Track, error)
	GetList(ctx context.Context, filters entities.TrackGetListFilters) ([]entities.Track, error)
}

type TracksHandlers struct {
	trackService TracksService
}

func NewTracksHandlers(g *echo.Group, ts TracksService) *TracksHandlers {
	ctl := &TracksHandlers{trackService: ts}
	g.POST("/", ctl.Create)
	g.GET("/", ctl.List)
	g.GET("/:id", ctl.Retrieve)
	return ctl
}

type HTTPError struct {
	Message string `json:"message"`
}
