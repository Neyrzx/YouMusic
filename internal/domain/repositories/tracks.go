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

func (r *TracksRepository) DeleteLyricByTrackID(ctx context.Context, tx pgx.Tx, trackID int) (err error) {
	sql := `DELETE FROM lyrics WHERE track_id = $1;`

	err = tx.QueryRow(ctx, sql, trackID).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to tx.QueryRow: %w", err)
	}

	return nil
}

func (r *TracksRepository) CreateLyric(ctx context.Context, tx pgx.Tx, lyrics []dao.Lyric) (err error) {
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

func (r *TracksRepository) CreateLyricFromSlice(ctx context.Context, tx pgx.Tx, trackID int, lyrics []string) (err error) {
	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"lyrics"},
		[]string{"track_id", "verse_text"},
		pgx.CopyFromSlice(len(lyrics), func(i int) ([]any, error) {
			return []any{trackID, lyrics[i]}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to insert lyrics: %w", err)
	}

	return nil
}

func (r *TracksRepository) CreateArtist(ctx context.Context, tx pgx.Tx, artist dao.Artist) (id int, err error) {
	sql := `
		INSERT INTO artists (name) VALUES ($1) RETURNING artist_id;`

	if err = tx.QueryRow(ctx, sql, artist.Name).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to rows.Scan: %w", err)
	}

	return id, nil
}

func (r *TracksRepository) CreateTrack(ctx context.Context, tx pgx.Tx, track dao.Track) (id int, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	sql := `
		INSERT INTO tracks (title, artist_id, link, released_at) VALUES ($1, $2, $3, $4) RETURNING track_id;`

	if err = tx.QueryRow(ctx, sql, track.Title, track.ArtistID, track.Link, track.ReleasedAt).Scan(&id); err != nil {
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

	sql := `
		SELECT EXISTS(
			SELECT 
				1
			FROM
				tracks JOIN artists ON tracks.artist_id = artists.artist_id
			WHERE
				tracks.title = $1 AND artists.name = $2
		);`

	if err = r.db.QueryRow(ctx, sql, track, artist).Scan(&exists); err != nil {
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

	sql := `
		SELECT artist_id FROM artists WHERE name = $1;`

	if err := tx.QueryRow(ctx, sql, name).Scan(&id); err != nil {
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

func (r *TracksRepository) GetArtistIDByTrackID(ctx context.Context, tx pgx.Tx, id int) (artistID int, err error) {
	sql := `SELECT artist_id FROM tracks WHERE track_id = $1;`

	if err = tx.QueryRow(ctx, sql, id).Scan(&artistID); err != nil {
		return 0, fmt.Errorf("failed to tx.QueryRow: %w", err)
	}

	return artistID, nil
}

func (r *TracksRepository) GetTrack(ctx context.Context, tx pgx.Tx, id int) (track entities.Track, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	sql := `
		SELECT artists.name, tracks.title, tracks.link, tracks.released_at
		FROM
			tracks JOIN artists
				ON tracks.artist_id = artists.artist_id
		WHERE
			tracks.track_id = $1;`

	if err = tx.QueryRow(ctx, sql, id).Scan(
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

	sql := `SELECT verse_text FROM lyrics WHERE track_id = $1;`

	rows, err := tx.Query(ctx, sql, trackID)
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

	sql := `SELECT track_id, verse_text FROM lyrics WHERE track_id in (`
	sql += strings.Join(params, ",")
	sql += ");"

	rows, err := tx.Query(ctx, sql, args...)
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

func (r *TracksRepository) UpdateArtist(ctx context.Context, tx pgx.Tx, artist dao.Artist) (err error) {
	sql := `UPDATE artists SET name = $1 WHERE artist_id = $2;`

	err = tx.QueryRow(ctx, sql, artist.Name, artist.ArtistID).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to tx.QueryRow: %w", err)
	}

	return nil
}

func (r *TracksRepository) UpdateTrack(ctx context.Context, tx pgx.Tx, track dao.Track) (err error) {
	var (
		sqlBase strings.Builder
		fields  []string
		args    = []any{track.TrackID}
		phIndex = 1
	)

	sqlBase.WriteString(`UPDATE tracks SET %s WHERE track_id = $1;`)

	if track.Title != "" {
		phIndex++
		fields = append(fields, fmt.Sprintf("title = $%d", phIndex))
		args = append(args, track.Title)
	}
	if track.Link != "" {
		phIndex++
		fields = append(fields, fmt.Sprintf("link = $%d", phIndex))
		args = append(args, track.Link)
	}
	if !track.ReleasedAt.IsZero() {
		phIndex++
		fields = append(fields, fmt.Sprintf("released_at = $%d", phIndex))
		args = append(args, track.ReleasedAt)
	}
	if len(fields) == 0 {
		return errors.New("failed to build sql empty track")
	}

	sql := fmt.Sprintf(sqlBase.String(), strings.Join(fields, ","))

	err = tx.QueryRow(ctx, sql, args...).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to tx.QueryRow: %w", err)
	}

	return nil
}

func (r *TracksRepository) DeleteTrackByID(ctx context.Context, trackID int) (err error) {
	sql := `DELETE FROM tracks WHERE track_id = $1;`

	err = r.db.QueryRow(ctx, sql, trackID).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to tx.QueryRow: %w", err)
	}

	return nil
}

func (r *TracksRepository) GetLyricPaginated(ctx context.Context, _ pgx.Tx, trackID int, offset int) (lyric dao.Lyric, err error) {
	sql := `SELECT lyric_id, verse_text FROM lyrics WHERE track_id = $1 LIMIT 1 OFFSET $2;`

	if err = r.db.QueryRow(ctx, sql, trackID, offset).Scan(
		&lyric.LyricID,
		&lyric.Verse,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dao.Lyric{}, domain.ErrTrackLyricNotFound
		}
		return dao.Lyric{}, err
	}

	return lyric, nil
}
