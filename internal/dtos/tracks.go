package dtos

import "time"

type TrackCreateDTO struct {
	Title  string
	Artist string
}

type TrackInfoDTO struct {
	Group string
	Song  string
}

type TrackInfoResultDTO struct {
	ReleaseDate time.Time
	Text        string
	Link        string
}
