package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	"github.com/neyrzx/youmusic/pkg/utils"
)

type TrackUpdateRequest struct {
	ID       int               `json:"-" param:"id"`
	Artist   string            `json:"artist"`
	Track    string            `json:"track"`
	Released utils.ReleaseDate `json:"released" format:"date" example:"10.10.2010"`
	Link     string            `json:"link" validate:"omitempty,uri" format:"uri" example:"https://y.be/asd2d2cW"`
	Lyric    string            `json:"lyric" example:"verse #1\n\nverse #2\n\nverse #3"`
}

// Update godoc
// @Summary      Update the tracks
// @Description  Updating the track
// @Tags         Tracks
// @Accept       json
// @Produce			 json
// @Param				 id path int true "track id"
// @Param				 input body v1.TrackUpdateRequest true "track id"
// @Success      204  {string}  string "OK"
// @Failure      400  {object}  v1.HTTPError "Bad request"
// @Failure      500  {object}  v1.HTTPError "Internal server error"
// @Router       /tracks/{id}/ [patch]
func (h *TracksHandlers) Update(c echo.Context) (err error) {
	var request TrackUpdateRequest

	if err = c.Bind(&request); err != nil {
		h.logger.Err(err).Msg("failed to c.Bind")
		return c.JSON(http.StatusBadRequest, HTTPError{Message: "request body malformed"})
	}

	if err = c.Validate(request); err != nil {
		h.logger.Err(err).Msg("failed to Validate")
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.trackService.Update(c.Request().Context(), entities.TrackUpdate{
		TrackID:  request.ID,
		Track:    request.Track,
		Artist:   request.Artist,
		Lyric:    request.Lyric,
		Link:     request.Link,
		Released: time.Time(request.Released),
	})
	if err != nil {
		h.logger.Err(err).Msg("failed to trackService.Update")
		return c.JSON(http.StatusInternalServerError, HTTPError{Message: "something went wrong"})
	}

	return c.JSON(http.StatusNoContent, "OK")
}
