package dao

import "time"

type ArtistDAO struct {
	ArtistID  int       `db:"artist_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type TrackDAO struct {
	TrackID    int       `db:"track_id"`
	ArtistID   int       `db:"artist_id"`
	Title      string    `db:"title"`
	Link       string    `db:"link"`
	ReleasedAt time.Time `db:"released_at"`
	CreatedAt  time.Time `db:"created_at"`
}

type LyricDAO struct {
	LyricID   int       `db:"lyric_id"`
	TrackID   int       `db:"track_id"`
	Verse     string    `db:"verse_text"`
	CreatedAt time.Time `db:"created_at"`
}
