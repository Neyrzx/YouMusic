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
	Released utils.ReleaseDate `json:"released"`
	Link     string            `json:"link" validate:"uri"`
	Lyric    string            `json:"lyric"`
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
		c.Logger().Errorf("failed to Bind: %s", err.Error())
		return c.JSON(http.StatusBadRequest, HTTPError{Message: "request body malformed"})
	}

	if err = c.Validate(request); err != nil {
		c.Logger().Errorf("failed to Validate: %s", err.Error())
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
		c.Logger().Errorf("failed to trackService.Update: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, HTTPError{Message: "something went wrong"})
	}

	return c.JSON(http.StatusNoContent, "OK")
}
