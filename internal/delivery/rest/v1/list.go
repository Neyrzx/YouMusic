package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neyrzx/youmusic/internal/domain/entities"
)

type TracksPaginationQuery struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type TracksListQuery struct {
	TracksPaginationQuery

	Artist       string `query:"artist"`
	Track        string `query:"track"`
	ReleasedYear string `query:"releasedyear"`
	Link         string `query:"link"`
}

type TracksResponse struct {
	Artist   string    `json:"artist"`
	Track    string    `json:"track"`
	Lyric    []string  `json:"lyric"`
	Link     string    `json:"link"`
	Released time.Time `json:"released"`
}

// List godoc
// @Summary      List of tracks
// @Description  List of tracks with filters
// @Tags         Tracks
// @Accept       json
// @Produce			 json
// @Param				 limit query string false "Limit result."
// @Param				 offset query string false "Offset result."
// @Param				 artist query string false "Name of the artist or group."
// @Param				 track query string false "Title of track."
// @Param				 releasedyaer query string false "List of tracks."
// @Param				 link query string false "Exact link"
// @Success      200  {array}  v1.TrackResponse "Success response"
// @Failure      400  {object}  v1.HTTPError "Bad request"
// @Failure      500  {object}  v1.HTTPError "Internal server error"
// @Router       /tracks/ [get]
func (h *TracksHandlers) List(c echo.Context) (err error) {
	var queryparam TracksListQuery

	if err = c.Bind(&queryparam); err != nil {
		return c.JSON(http.StatusBadRequest, HTTPError{Message: err.Error()})
	}

	if err = c.Validate(queryparam); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var tracks []entities.Track
	if tracks, err = h.trackService.GetList(c.Request().Context(), entities.TrackGetListFilters{
		Limit:        queryparam.Limit,
		Offset:       queryparam.Offset,
		Artist:       queryparam.Artist,
		Track:        queryparam.Track,
		ReleasedYear: queryparam.ReleasedYear,
		Link:         queryparam.Link,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, HTTPError{Message: "something went wrong"})
	}

	res := []TracksResponse{}
	for _, track := range tracks {
		res = append(res, TracksResponse{
			Artist:   track.Artist,
			Track:    track.Track,
			Lyric:    track.Lyric,
			Link:     track.Link,
			Released: track.Released,
		})
	}

	return c.JSON(http.StatusOK, res)
}
