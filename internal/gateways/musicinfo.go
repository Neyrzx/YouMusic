package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/neyrzx/youmusic/internal/config"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	"github.com/neyrzx/youmusic/pkg/utils"
)

const timeoutInfo = 15 * time.Second

type Client interface {
	Get(context.Context, string) (io.ReadCloser, error)
}

type MusicInfoGateway struct {
	client Client
	cfg    config.GatewayHTTPClient
}

func NewMusicInfoGateway(client Client, cfg config.GatewayHTTPClient) *MusicInfoGateway {
	return &MusicInfoGateway{client: client, cfg: cfg}
}

type InfoResponse struct {
	ReleaseDate utils.ReleaseDate `json:"releaseDate"`
	Text        string            `json:"text"`
	Link        string            `json:"link"`
}

func (gw *MusicInfoGateway) Info(ctx context.Context, track entities.TrackInfo) (entities.TrackInfoResult, error) {
	_, cancelFunc := context.WithTimeout(ctx, timeoutInfo)
	defer cancelFunc()

	path := fmt.Sprintf("%s/info/?song=%s&group=%s", gw.cfg.BaseURL, track.Song, track.Group)
	data, err := gw.client.Get(ctx, path)
	if err != nil {
		return entities.TrackInfoResult{}, fmt.Errorf("failed to client.Get(%s): %w", path, err)
	}
	defer data.Close()

	var response InfoResponse
	if err = json.NewDecoder(data).Decode(&response); err != nil {
		return entities.TrackInfoResult{}, fmt.Errorf("failed to json.Decode(): %w", err)
	}

	return entities.TrackInfoResult{
		ReleaseDate: time.Time(response.ReleaseDate),
		Text:        response.Text,
		Link:        response.Link,
	}, nil
}
