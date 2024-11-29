package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/neyrzx/youmusic/internal/config"
	"github.com/neyrzx/youmusic/internal/dtos"
	"github.com/neyrzx/youmusic/pkg/utils"
)

const timeoutInfo = 15 * time.Second

type Client interface {
	Get(string) (io.ReadCloser, error)
}

type MusicInfoGateway struct {
	client Client
	cfg    config.MusicInfoGateway
}

func NewMusicInfoGateway(client Client, cfg config.MusicInfoGateway) *MusicInfoGateway {
	return &MusicInfoGateway{client: client, cfg: cfg}
}

type InfoResponse struct {
	ReleaseDate utils.ReleaseDate `json:"releaseDate"`
	Text        string            `json:"text"`
	Link        string            `json:"link"`
}

func (gw *MusicInfoGateway) Info(ctx context.Context, track dtos.TrackInfoDTO) (dtos.TrackInfoResultDTO, error) {
	_, cancelFunc := context.WithTimeout(ctx, timeoutInfo)
	defer cancelFunc()

	path := fmt.Sprintf("%s%s?song=%s&group=%s", gw.cfg.URL, gw.cfg.Route, track.Song, track.Group)
	data, err := gw.client.Get(path)
	if err != nil {
		return dtos.TrackInfoResultDTO{}, fmt.Errorf("failed to client.Get(%s): %w", path, err)
	}
	defer data.Close()

	var response InfoResponse
	if err = json.NewDecoder(data).Decode(&response); err != nil {
		return dtos.TrackInfoResultDTO{}, fmt.Errorf("failed to json.Decode(): %w", err)
	}

	return dtos.TrackInfoResultDTO{
		ReleaseDate: time.Time(response.ReleaseDate),
		Text:        response.Text,
		Link:        response.Link,
	}, nil
}
