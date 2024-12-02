package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TrackDeletePrarams struct {
	ID int `json:"-" param:"id"`
}

// Delete godoc
// @Summary      Delete track
// @Description  Deliting track by track id
// @Tags         Tracks
// @Accept       json
// @Produce			 json
// @Param				 id path int true "track id"
// @Success      204  {string} string "OK"
// @Failure      400  {object}  v1.HTTPError "Bad request"
// @Failure      500  {object}  v1.HTTPError "Internal server error"
// @Router       /tracks/{id}/ [delete]
func (h *TracksHandlers) Delete(c echo.Context) (err error) {
	var request TrackDeletePrarams

	if err = c.Bind(&request); err != nil {
		h.logger.Err(err).Msg("failed to c.Bind")
		return c.JSON(http.StatusBadRequest, HTTPError{Message: "request body malformed"})
	}

	if err = c.Validate(request); err != nil {
		h.logger.Err(err).Msg("failed to c.Validate")
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = h.trackService.Delete(c.Request().Context(), request.ID); err != nil {
		h.logger.Err(err).Msg("failed to trackService.Update")
		return c.JSON(http.StatusInternalServerError, HTTPError{Message: "something went wrong"})
	}

	return c.JSON(http.StatusNoContent, "OK")
}
