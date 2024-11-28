package services

import (
	"context"
	"fmt"
	"time"

	"github.com/neyrzx/youmusic/internal/domain/entities"
	"github.com/neyrzx/youmusic/internal/dtos"
	"github.com/neyrzx/youmusic/pkg/utils"
)

const createMethodTimeout = 15 * time.Second

type TracksRepository interface {
	Create(ctx context.Context, track entities.Track) error
	Exists(ctx context.Context, track entities.Track) error
}

type TracksInfoGateway interface {
	Info(ctx context.Context, track dtos.TrackInfoDTO) (dtos.TrackInfoResultDTO, error)
}

type TracksService struct {
	repo        TracksRepository
	infoGateway TracksInfoGateway
}

func NewTracksService(repo TracksRepository, infoGateway TracksInfoGateway) *TracksService {
	return &TracksService{repo: repo, infoGateway: infoGateway}
}

func (s *TracksService) Create(ctx context.Context, track dtos.TrackCreateDTO) (err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, createMethodTimeout)
	defer cancelFunc()

	if err = s.repo.Exists(ctx, entities.Track{
		Title:  track.Title,
		Artist: track.Artist,
	}); err != nil {
		return fmt.Errorf("failed to repo.Exists(): %w", err)
	}

	var trackInfo dtos.TrackInfoResultDTO
	if trackInfo, err = s.infoGateway.Info(ctx, dtos.TrackInfoDTO{
		Group: track.Artist,
		Song:  track.Title,
	}); err != nil {
		return fmt.Errorf("failed to infoGateway.Info(): %w", err)
	}

	if err = s.repo.Create(ctx, entities.Track{
		Title:    track.Title,
		Artist:   track.Artist,
		Lyrics:   utils.SplitLyricsToVerses(ctx, trackInfo.Text),
		Link:     trackInfo.Link,
		Released: trackInfo.ReleaseDate,
	}); err != nil {
		return fmt.Errorf("failed to repo.Create(): %w", err)
	}

	return nil
}
