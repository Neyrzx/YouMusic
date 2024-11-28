package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/neyrzx/youmusic/internal/dtos"
	"github.com/neyrzx/youmusic/pkg/utils"
)

const timeoutInfo = 15 * time.Second

type Client interface {
	Get(string) (*http.Response, error)
}

type MusicInfoGateway struct {
	client Client
}

func NewMusicInfoGateway(client Client) *MusicInfoGateway {
	return &MusicInfoGateway{client: client}
}

type InfoResponse struct {
	ReleaseDate utils.ReleaseDate `json:"releaseDate"`
	Text        string            `json:"text"`
	Link        string            `json:"link"`
}

func (gw *MusicInfoGateway) Info(ctx context.Context, track dtos.TrackInfoDTO) (*dtos.TrackInfoResultDTO, error) {
	_, cancelFunc := context.WithTimeout(ctx, timeoutInfo)
	defer cancelFunc()

	path := fmt.Sprintf("/info?song=%s&group=%s", track.Song, track.Group)
	res, err := gw.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to client.Get(%s): %w", path, err)
	}
	defer res.Body.Close()

	var response InfoResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to json.Decode(): %w", err)
	}

	return &dtos.TrackInfoResultDTO{
		ReleaseDate: time.Time(response.ReleaseDate),
		Text:        response.Text,
		Link:        response.Link,
	}, nil
}
