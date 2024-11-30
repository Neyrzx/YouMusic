package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	domain "github.com/neyrzx/youmusic/internal/domain/errors"
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

func (r *TracksRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin trasaction r.db.Begin: %w", err)
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				//TODO: залогировать
				// fmt.Errorf("failed to tx.Rollback: %w", err)
				_ = txErr
			}
		}
	}()

	if err = fn(tx); err != nil {
		return fmt.Errorf("failed while execute transaction: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *TracksRepository) CreateLyric(ctx context.Context, tx pgx.Tx, lyrics []dao.Lyric) (err error) {
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

func (r *TracksRepository) CreateArtist(ctx context.Context, tx pgx.Tx, artist dao.Artist) (id int, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		INSERT INTO artists (name) VALUES ($1) RETURNING artist_id;`

	if err = tx.QueryRow(ctx, query, artist.Name).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to rows.Scan: %w", err)
	}

	return id, nil
}

func (r *TracksRepository) CreateTrack(ctx context.Context, tx pgx.Tx, track dao.Track) (id int, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		INSERT INTO tracks (title, artist_id, link, released_at) VALUES ($1, $2, $3, $4) RETURNING track_id;`

	if err = tx.QueryRow(ctx, query, track.Title, track.ArtistID, track.Link, track.ReleasedAt).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "tracks_title_artist_id_unique" {
			return 0, domain.ErrTrackAlreadyExists
		}
		return 0, fmt.Errorf("failed to rows.Scan: %w", err)
	}

	return id, nil
}

func (r *TracksRepository) IsTrackExists(ctx context.Context, track string, artist string) (exists bool, err error) {
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

	if err = r.db.QueryRow(ctx, query, track, artist).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}

		return false, fmt.Errorf("failed to rows.Scan: %w", err)
	}

	return exists, nil
}

func (r *TracksRepository) IsArtistExists(ctx context.Context, tx pgx.Tx, name string) (id int, exists bool) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		SELECT artist_id FROM artists WHERE name = $1;`

	if err := tx.QueryRow(ctx, query, name).Scan(&id); err != nil {
		return 0, false
	}

	return id, true
}

func (r *TracksRepository) GetByID(ctx context.Context, id int) (track entities.Track, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entities.Track{}, fmt.Errorf("failed to begin trasaction r.db.Begin: %w", err)
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				//TODO: залогировать
				// fmt.Errorf("failed to tx.Rollback: %w", err)
				_ = txErr
			}
		}
	}()

	if track, err = r.GetTrack(ctx, tx, id); err != nil {
		return entities.Track{}, fmt.Errorf("failed to GetTrack(%d): %w", id, err)
	}

	if track.Lyric, err = r.GetTrackLyric(ctx, tx, id); err != nil {
		return entities.Track{}, fmt.Errorf("failed to GetTrackLyric(%d): %w", id, err)
	}

	return track, tx.Commit(ctx)
}

func (r *TracksRepository) GetTrack(ctx context.Context, tx pgx.Tx, id int) (track entities.Track, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `
		SELECT artists.name, tracks.title, tracks.link, tracks.released_at
		FROM
			tracks JOIN artists
				ON tracks.artist_id = artists.artist_id
		WHERE
			tracks.track_id = $1;`

	if err = tx.QueryRow(ctx, query, id).Scan(
		&track.Artist,
		&track.Track,
		&track.Link,
		&track.Released,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Track{}, domain.ErrTrackNotFound
		}

		return entities.Track{}, fmt.Errorf("failed to QueryRow(%d): %w", id, err)
	}

	return track, nil
}

func (r *TracksRepository) GetTracksByFilter(ctx context.Context, tx pgx.Tx, filter entities.TrackGetListFilters) (tracks map[int]entities.Track, err error) {
	var (
		sqlBase      strings.Builder
		clause       []string
		args         []any
		paramIdx     = 1
		defaultLimit = 10
	)

	sqlBase.WriteString(`
	SELECT
		tracks.track_id,
		artists.name,
		tracks.title,
		tracks.released_at,
		tracks.link
	FROM
		tracks JOIN artists ON tracks.artist_id = artists.artist_id
	`)

	if filter.Artist != "" {
		clause = append(clause, fmt.Sprintf(`artists.name = $%d`, paramIdx))
		paramIdx++
		args = append(args, filter.Artist)
	}

	if filter.Track != "" {
		clause = append(clause, fmt.Sprintf(`tracks.title = $%d`, paramIdx))
		paramIdx++
		args = append(args, filter.Track)
	}

	if filter.Link != "" {
		clause = append(clause, fmt.Sprintf(`tracks.link = $%d`, paramIdx))
		paramIdx++
		args = append(args, filter.Link)
	}

	if filter.ReleasedYear != "" {
		clause = append(clause, fmt.Sprintf(`EXTRACT(YEAR FROM tracks.released_at) = $%d`, paramIdx))
		paramIdx++
		args = append(args, filter.ReleasedYear)
	}

	if len(clause) > 0 {
		sqlBase.WriteString(`WHERE `)
		sqlBase.WriteString(strings.Join(clause, " AND "))
		sqlBase.WriteString(` `)
	}

	sqlBase.WriteString(`ORDER BY tracks.track_id ASC `)

	switch {
	case filter.Limit == 0 && filter.Offset == 0:
		sqlBase.WriteString(fmt.Sprintf(` LIMIT $%d `, paramIdx))
		args = append(args, defaultLimit)
	case filter.Limit != 0 && filter.Offset == 0:
		sqlBase.WriteString(fmt.Sprintf(` LIMIT $%d `, paramIdx))
		args = append(args, filter.Limit)
	case filter.Limit != 0 && filter.Offset != 0:
		sqlBase.WriteString(fmt.Sprintf(` LIMIT $%d `, paramIdx))
		paramIdx++
		sqlBase.WriteString(fmt.Sprintf(` OFFSET $%d `, paramIdx))
		args = append(args, filter.Limit, filter.Offset)
	}

	sqlBase.WriteString(`;`)
	sql := sqlBase.String()

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to tracks tx.Query: %w", err)
	}

	tracks = make(map[int]entities.Track)
	track := entities.Track{}

	for rows.Next() {
		if err = rows.Scan(
			&track.ID,
			&track.Artist,
			&track.Track,
			&track.Released,
			&track.Link,
		); err != nil {
			return nil, fmt.Errorf("failed to rows.Scan: %w", err)
		}
		tracks[track.ID] = track
	}

	return tracks, nil
}

func (r *TracksRepository) GetTrackLyric(ctx context.Context, tx pgx.Tx, trackID int) (lyric []string, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	query := `SELECT verse_text FROM lyrics WHERE track_id = $1;`

	rows, err := tx.Query(ctx, query, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to Query(%d): %w", trackID, err)
	}

	var verse string
	for rows.Next() {
		if err = rows.Scan(&verse); err != nil {
			return nil, fmt.Errorf("failed while scanning query result: %w", err)
		}
		lyric = append(lyric, verse)
	}

	return lyric, nil
}

func (r *TracksRepository) GetLyricsByTrackIDs(ctx context.Context, tx pgx.Tx, ids []int) (lyrics []dao.Lyric, err error) {
	var (
		params []string
		args   []any
	)

	if len(ids) == 0 {
		return nil, nil
	}

	for i, id := range ids {
		params = append(params, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := `SELECT track_id, verse_text FROM lyrics WHERE track_id in (`
	query += strings.Join(params, ",")
	query += ");"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to tx.Query: %w", err)
	}

	var lyric dao.Lyric
	for rows.Next() {
		if err = rows.Scan(&lyric.TrackID, &lyric.Verse); err != nil {
			return nil, fmt.Errorf("failed to rows.Scan: %w", err)
		}
		lyrics = append(lyrics, lyric)
	}

	return lyrics, nil
}
