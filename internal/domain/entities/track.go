package entities

import "time"

type Track struct {
	Title    string
	Artist   string
	Lyrics   []string
	Link     string
	Released time.Time
}
