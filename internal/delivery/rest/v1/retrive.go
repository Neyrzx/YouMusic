package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	domain "github.com/neyrzx/youmusic/internal/domain/errors"
)

type TracksRetrievePathParam struct {
	ID int `param:"id"`
}

type TracksRetrieveResponse struct {
	Artist   string    `json:"artist"`
	Track    string    `json:"track"`
	Lyric    []string  `json:"lyric"`
	Link     string    `json:"link"`
	Released time.Time `json:"released"`
}

// Retrieve godoc
// @Summary      Retrive track
// @Description  Retriving track
// @Tags         Tracks
// @Accept       json
// @Produce			 json
// @Param				 id path int true "track id"
// @Success      200  {object}  v1.TracksRetrieveResponse "get track result"
// @Failure      400  {object}  v1.HTTPError "Bad request"
// @Failure      422  {object}  v1.HTTPError "Validation errors"
// @Failure      500  {object}  v1.HTTPError "Internal server error"
// @Router       /tracks/{id}/ [get]
func (h *TracksHandlers) Retrieve(c echo.Context) (err error) {
	var pathParam TracksRetrievePathParam

	if err = c.Bind(&pathParam); err != nil {
		h.logger.Err(err).Msg("failed to c.Bind")
		return c.JSON(http.StatusBadRequest, HTTPError{Message: "id param is invalid"})
	}

	if err = c.Validate(pathParam); err != nil {
		h.logger.Err(err).Msg("failed to c.Validate")
		return c.JSON(http.StatusBadRequest, err)
	}

	track, err := h.trackService.GetByID(c.Request().Context(), pathParam.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTrackNotFound) {
			return c.JSON(http.StatusNotFound, err)
		}
		h.logger.Err(err).Int("trackID", pathParam.ID).Msg("failed to c.Validate")
		return c.JSON(http.StatusInternalServerError, HTTPError{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, TracksRetrieveResponse{
		Artist:   track.Artist,
		Track:    track.Track,
		Lyric:    track.Lyric,
		Link:     track.Link,
		Released: track.Released,
	})
}
