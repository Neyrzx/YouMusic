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
// @Success      201  {string}  string "Success created"
// @Failure      400  {object}  v1.HTTPError "Bad request"
// @Failure      500  {object}  v1.HTTPError "Internal server error"
// @Router       /tracks/ [post]
func (h *TracksHandlers) Create(c echo.Context) (err error) {
	var request TracksCreateRequest

	if err = c.Bind(&request); err != nil {
		h.logger.Err(err).Msg("failed to c.Bind")
		return c.JSON(http.StatusBadRequest, HTTPError{Message: "request body malformed"})
	}

	if err = c.Validate(request); err != nil {
		h.logger.Err(err).Msg("failed to c.Validate")
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = h.trackService.Create(c.Request().Context(), entities.TrackCreate{
		Title:  request.Song,
		Artist: request.Group,
	}); err != nil {
		h.logger.Err(err).Msg("failed to trackService.Create")
		return c.JSON(http.StatusInternalServerError, HTTPError{Message: "something went wrong, try again later"})
	}

	return c.JSON(http.StatusCreated, "OK")
}
