package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	"github.com/neyrzx/youmusic/internal/domain/repositories/dao"
)

const queryTimeout = 120 * time.Second

var (
	ErrTrackAlreadyExists = errors.New("track is already exists")
)

type TracksRepository struct {
	db *pgxpool.Pool
}

func NewTracksRepository(db *pgxpool.Pool) *TracksRepository {
	return &TracksRepository{db: db}
}

func (r *TracksRepository) Create(ctx context.Context, track entities.Track) (err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin trasaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	trackExists, err := r.IsExistsTrack(ctx, tx, track)
	if err != nil {
		return fmt.Errorf("failed to GetTrack(%s, %s): %w", track.Title, track.Artist, err)
	}

	if trackExists {
		return ErrTrackAlreadyExists
	}

	trackDAO := dao.TrackDAO{
		Title:      track.Title,
		Link:       track.Link,
		ReleasedAt: track.Released,
	}
	artistID, trackExists := r.IsExistsArtist(ctx, tx, track.Artist)
	if trackExists {
		trackDAO.ArtistID = artistID
	} else {
		trackDAO.ArtistID, err = r.CreateArtist(ctx, tx, dao.ArtistDAO{Name: track.Artist})
		if err != nil {
			return fmt.Errorf("failed to CreateArtist(%s): %w", track.Artist, err)
		}
	}

	trackID, err := r.CreateTrack(ctx, tx, trackDAO)
	if err != nil {
		return fmt.Errorf("failed to CreateTrack(%v+): %w", trackDAO, err)
	}

	lyricsDAO := make([]dao.LyricDAO, len(track.Lyrics))
	for i, verse := range track.Lyrics {
		lyricsDAO[i].TrackID = trackID
		lyricsDAO[i].Verse = verse
	}
	if err = r.CreateLyric(ctx, tx, lyricsDAO); err != nil {
		return fmt.Errorf("failed to CreateLyric for artist (%d, %s): %w", artistID, track.Artist, err)
	}

	return tx.Commit(ctx)
}

func (r *TracksRepository) CreateLyric(ctx context.Context, tx pgx.Tx, lyrics []dao.LyricDAO) (err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"lyrics"},
		[]string{"track_id", "verse_text"},
		pgx.CopyFromSlice(len(lyrics), func(i int) ([]any, error) {
			return []any{lyrics[i].TrackID, lyrics[i].Verse}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to insert lyrics: %w", err)
	}

	return nil
}

func (r *TracksRepository) CreateArtist(ctx context.Context, tx pgx.Tx, artist dao.ArtistDAO) (id int, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		INSERT INTO artists (name) VALUES ($1) RETURNING artist_id;`

	if err = tx.QueryRow(ctx, query, artist.Name).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to rows.Scan: %w", err)
	}

	return id, nil
}

func (r *TracksRepository) CreateTrack(ctx context.Context, tx pgx.Tx, track dao.TrackDAO) (id int, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		INSERT INTO tracks (title, artist_id, link, released_at) VALUES ($1, $2, $3, $4) RETURNING track_id;`

	if err = tx.QueryRow(ctx, query, track.Title, track.ArtistID, track.Link, track.ReleasedAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to rows.Scan: %w", err)
	}

	return id, nil
}

func (r *TracksRepository) Exists(_ context.Context, _ entities.Track) error {
	return nil
}

func (r *TracksRepository) IsExistsTrack(
	ctx context.Context, tx pgx.Tx, track entities.Track,
) (exists bool, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		SELECT EXISTS(
			SELECT 
				1
			FROM
				tracks JOIN artists ON tracks.artist_id = artists.artist_id
			WHERE
				tracks.title = $1 AND artists.name = $2
		);`

	if err = tx.QueryRow(ctx, query, track.Title, track.Artist).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}

		return false, fmt.Errorf("failed to rows.Scan: %w", err)
	}

	return exists, nil
}

func (r *TracksRepository) IsExistsArtist(ctx context.Context, tx pgx.Tx, name string) (id int, exists bool) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		SELECT artist_id FROM artists WHERE name = $1;`

	if err := tx.QueryRow(ctx, query, name).Scan(&id); err != nil {
		return 0, false
	}

	return id, true
}
