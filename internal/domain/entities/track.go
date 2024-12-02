package entities

import "time"

type Track struct {
	ID       int
	Track    string
	Artist   string
	Lyric    []string
	Link     string
	Released time.Time
}

type TrackVerse struct {
	OrderID int
	Verse   string
}

type TrackCreate struct {
	Title  string
	Artist string
}

type TrackUpdate struct {
	TrackID  int
	Track    string
	Artist   string
	Lyric    string
	Link     string
	Released time.Time
}

type TrackGetListFilters struct {
	Limit        int
	Offset       int
	Artist       string
	Track        string
	ReleasedYear string
	Link         string
}

type TrackInfoResult struct {
	ReleaseDate time.Time
	Text        string
	Link        string
}

type TrackInfo struct {
	Group string
	Song  string
}
