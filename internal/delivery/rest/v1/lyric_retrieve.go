package v1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	domain "github.com/neyrzx/youmusic/internal/domain/errors"
)

type TrackLyricParams struct {
	TrackID int `param:"id"`
	Offset  int `query:"offset"`
}

type TrackLyricResponse struct {
	OrderID int    `json:"orderID"`
	Verse   string `json:"verse"`
}

// LyricRetrieve godoc
// @Summary      Retrive verse
// @Description  Retrive lyric verse with offset
// @Tags         Tracks
// @Accept       json
// @Produce			 json
// @Param				 id path int true "track id"
// @Param				 offset query int false "verse offset"
// @Success      200  {object}  v1.TrackLyricResponse "Success response"
// @Failure      400  {object}  v1.HTTPError "Bad request"
// @Failure      422  {object}  v1.HTTPError "Validation errors"
// @Failure      500  {object}  v1.HTTPError "Internal server error"
// @Router       /tracks/{id}/lyric/ [get]
func (h *TracksHandlers) LyricRetrieve(c echo.Context) (err error) {
	var request TrackLyricParams

	if err = c.Bind(&request); err != nil {
		h.logger.Err(err).Msg("failed to c.Bind")
		return c.JSON(http.StatusBadRequest, HTTPError{Message: "request body malformed"})
	}

	if err = c.Validate(request); err != nil {
		h.logger.Err(err).Msg("failed to c.Validate")
		return c.JSON(http.StatusBadRequest, err)
	}

	verse, err := h.trackService.GetLyric(c.Request().Context(), request.TrackID, request.Offset)
	if err != nil {
		h.logger.Err(err).
			Int("trackID", request.TrackID).
			Int("offset", request.Offset).
			Msg("failed to trackService.GetLyric")

		if errors.Is(err, domain.ErrTrackLyricNotFound) {
			return c.JSON(http.StatusNotFound, HTTPError{Message: domain.ErrTrackLyricNotFound.Error()})
		}

		return c.JSON(http.StatusInternalServerError, HTTPError{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, verse)
}
