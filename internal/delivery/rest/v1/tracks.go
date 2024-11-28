package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neyrzx/youmusic/internal/dtos"
)

type TracksService interface {
	Create(ctx context.Context, track dtos.TrackCreateDTO) error
}

type TracksHandlers struct {
	tc TracksService
}

func NewTracksHandlers(g *echo.Group, tc TracksService) *TracksHandlers {
	ctl := &TracksHandlers{tc: tc}
	g.POST("/", ctl.Create)
	return ctl
}

type TracksCreateRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

// TODO: Возвращать ошибки в более читаемом формате, т.к. сырые ошибки валидации имеют избыточные данные.
func (ctl *TracksHandlers) Create(c echo.Context) (err error) {
	var request TracksCreateRequest

	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "request body malformed"})
	}

	if err = c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	if err = ctl.tc.Create(c.Request().Context(), dtos.TrackCreateDTO{
		Title:  request.Song,
		Artist: request.Group,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}
