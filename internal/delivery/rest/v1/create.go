package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neyrzx/youmusic/internal/domain/entities"
)

type TracksCreateRequest struct {
	Group string `json:"group" validate:"required" example:"Muse"`
	Song  string `json:"song" validate:"required" example:"Song name"`
}

// Create godoc
// @Summary      Create track
// @Description  Creating track
// @Tags         Tracks
// @Accept       json
// @Produce			 json
// @Param				 input body v1.TracksCreateRequest true "Create track by song and group names."
// @Success      200  {string}  string "Success created"
// @Failure      400  {object}  v1.HTTPError "Bad request"
// @Failure      500  {object}  v1.HTTPError "Internal server error"
// @Router       /tracks/ [post]
func (h *TracksHandlers) Create(c echo.Context) (err error) {
	var request TracksCreateRequest

	if err = c.Bind(&request); err != nil {
		c.Logger().Errorf("failed to Create.Bind: %s", err.Error())
		return c.JSON(http.StatusBadRequest, HTTPError{Message: "request body malformed"})
	}

	if err = c.Validate(request); err != nil {
		c.Logger().Errorf("failed to Create.c.Validate: %s", err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = h.trackService.Create(c.Request().Context(), entities.TrackCreate{
		Title:  request.Song,
		Artist: request.Group,
	}); err != nil {
		c.Logger().Errorf("failed to Create.ctl.tc.Create: %s", err.Error())
		return c.JSON(http.StatusBadRequest, HTTPError{Message: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}
