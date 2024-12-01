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
	Update(ctx context.Context, track entities.TrackUpdate) error
	Delete(ctx context.Context, trackID int) error
	GetLyric(ctx context.Context, trackID int, offset int) (entities.TrackVerse, error)
}

type TracksHandlers struct {
	trackService TracksService
}

func NewTracksHandlers(g *echo.Group, ts TracksService) *TracksHandlers {
	h := &TracksHandlers{trackService: ts}
	g.POST("/", h.Create)
	g.GET("/", h.List)
	g.GET("/:id/", h.Retrieve)
	g.PATCH("/:id/", h.Update)
	g.DELETE("/:id/", h.Delete)
	g.GET("/:id/lyric/", h.LyricRetrieve)
	return h
}

type HTTPError struct {
	Message string `json:"message"`
}
