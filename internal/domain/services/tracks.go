package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	domain "github.com/neyrzx/youmusic/internal/domain/errors"
	"github.com/neyrzx/youmusic/internal/domain/repositories/dao"
	"github.com/neyrzx/youmusic/pkg/utils"
)

const createMethodTimeout = 15 * time.Second

type TracksRepository interface {
	CreateArtist(ctx context.Context, tx pgx.Tx, artist dao.Artist) (id int, err error)
	CreateTrack(ctx context.Context, tx pgx.Tx, track dao.Track) (id int, err error)
	CreateLyric(ctx context.Context, tx pgx.Tx, lyrics []dao.Lyric) (err error)
	GetByID(ctx context.Context, ID int) (entities.Track, error)
	GetTracksByFilter(ctx context.Context, tx pgx.Tx, filter entities.TrackGetListFilters) (tracks map[int]entities.Track, err error)
	GetLyricsByTrackIDs(ctx context.Context, tx pgx.Tx, IDs []int) (lyrics []dao.Lyric, err error)
	IsTrackExists(ctx context.Context, trackName string, artistName string) (exists bool, err error)
	IsArtistExists(ctx context.Context, tx pgx.Tx, name string) (id int, exists bool)
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type TracksInfoGateway interface {
	Info(ctx context.Context, track entities.TrackInfo) (entities.TrackInfoResult, error)
}

type TracksService struct {
	repo        TracksRepository
	infoGateway TracksInfoGateway
}

func NewTracksService(repo TracksRepository, infoGateway TracksInfoGateway) *TracksService {
	return &TracksService{repo: repo, infoGateway: infoGateway}
}

func (s *TracksService) Create(ctx context.Context, track entities.TrackCreate) (err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, createMethodTimeout)
	defer cancelFunc()

	trackExists, err := s.repo.IsTrackExists(ctx, track.Title, track.Artist)
	if err != nil {
		return fmt.Errorf("failed to repo.IsTrackExists(%s, %s): %w", track.Title, track.Artist, err)
	}
	if trackExists {
		return domain.ErrTrackAlreadyExists
	}

	var trackInfo entities.TrackInfoResult
	if trackInfo, err = s.infoGateway.Info(ctx, entities.TrackInfo{
		Group: track.Artist,
		Song:  track.Title,
	}); err != nil {
		return fmt.Errorf("failed to infoGateway.Info(): %w", err)
	}

	err = s.repo.WithTx(ctx, func(tx pgx.Tx) error {
		trackDAO := dao.Track{
			Title:      track.Title,
			Link:       trackInfo.Link,
			ReleasedAt: trackInfo.ReleaseDate,
		}
		artistID, exists := s.repo.IsArtistExists(ctx, tx, track.Artist)
		if exists {
			trackDAO.ArtistID = artistID
		} else {
			trackDAO.ArtistID, err = s.repo.CreateArtist(ctx, tx, dao.Artist{Name: track.Artist})
			if err != nil {
				return fmt.Errorf("failed to CreateArtist(%s): %w", track.Artist, err)
			}
		}

		var trackID int
		trackID, err = s.repo.CreateTrack(ctx, tx, trackDAO)
		if err != nil {
			return fmt.Errorf("failed to CreateTrack(%v+): %w", trackDAO, err)
		}

		lyric := utils.SplitLyricsToVerses(ctx, trackInfo.Text)
		lyricsDAO := make([]dao.Lyric, len(lyric))
		for i, verse := range lyric {
			lyricsDAO[i].TrackID = trackID
			lyricsDAO[i].Verse = verse
		}

		if err = s.repo.CreateLyric(ctx, tx, lyricsDAO); err != nil {
			return fmt.Errorf("failed to CreateLyric for artist (%d, %s): %w", artistID, track.Artist, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed from repo create track: %w", err)
	}

	return nil
}

func (s *TracksService) GetByID(ctx context.Context, id int) (track entities.Track, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, createMethodTimeout)
	defer cancelFunc()

	if track, err = s.repo.GetByID(ctx, id); err != nil {
		return entities.Track{}, fmt.Errorf("failed to repo.GetByID: %w", err)
	}

	return track, nil
}

func (s *TracksService) GetList(ctx context.Context, filters entities.TrackGetListFilters) (tracks []entities.Track, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, createMethodTimeout)
	defer cancelFunc()

	err = s.repo.WithTx(ctx, func(tx pgx.Tx) error {
		var tracksMap map[int]entities.Track
		tracksMap, err = s.repo.GetTracksByFilter(ctx, tx, filters)
		if err != nil {
			return fmt.Errorf("failed while getting tracks by filter: %w", err)
		}

		var IDs []int
		for ID := range tracksMap {
			IDs = append(IDs, ID)
		}

		var lyrics []dao.Lyric
		lyrics, err = s.repo.GetLyricsByTrackIDs(ctx, tx, IDs)
		if err != nil {
			return fmt.Errorf("failed while getting lyrics for tracks by IDs: %w", err)
		}

		var track entities.Track
		for _, lyric := range lyrics {
			track = tracksMap[lyric.TrackID]
			track.Lyric = append(track.Lyric, lyric.Verse)
			tracksMap[lyric.TrackID] = track
		}

		for _, track := range tracksMap {
			tracks = append(tracks, track)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to repo.GetByFilter: %w", err)
	}

	return tracks, nil
}
