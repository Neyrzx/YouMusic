package dao

import "time"

type Artist struct {
	ArtistID  int
	Name      string
	CreatedAt time.Time
}

type Track struct {
	TrackID    int
	ArtistID   int
	Title      string
	Link       string
	ReleasedAt time.Time
	CreatedAt  time.Time
}

type Lyric struct {
	LyricID   int
	TrackID   int
	Verse     string
	CreatedAt time.Time
}
