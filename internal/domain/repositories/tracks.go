package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	"github.com/neyrzx/youmusic/internal/domain/repositories/dao"
)

const queryTimeout = 120 * time.Second

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
			err = tx.Rollback(ctx)
		}
	}()

	artistdao := dao.ArtistDAO{Name: track.Artist}
	artistID, err := r.CreateArtist(ctx, tx, artistdao)
	if err != nil {
		return fmt.Errorf("failed to CreateArtist(%s): %w", track.Artist, err)
	}

	trackDAO := dao.TrackDAO{
		ArtistID:   artistID,
		Title:      track.Title,
		Link:       track.Link,
		ReleasedAt: track.Released,
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

	query := `INSERT INTO public.artists (name) VALUES ($1) RETURNING artist_id;`

	rows, err := tx.Query(ctx, query, artist.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to tx.Query: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("failed to rows.Scan: %w", err)
		}
	}

	return id, nil
}

func (r *TracksRepository) CreateTrack(ctx context.Context, tx pgx.Tx, track dao.TrackDAO) (id int, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `INSERT INTO public.tracks (title, artist_id, link, released_at) VALUES ($1, $2, $3, $4) RETURNING track_id;`

	rows, err := tx.Query(ctx, query, track.Title, track.ArtistID, track.Link, track.ReleasedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to tx.Query: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("failed to rows.Scan: %w", err)
		}
	}

	return id, nil
}

func (r *TracksRepository) Exists(_ context.Context, _ entities.Track) error {
	return nil
}
