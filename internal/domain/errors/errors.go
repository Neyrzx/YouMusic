package errors

import "errors"

var (
	ErrTrackAlreadyExists     = errors.New("track already exists")
	ErrTrackRequestInfoFailed = errors.New("failed to request the track info from external API")
	ErrTrackFailedCreateTrack = errors.New("failed to save the tack into DB")
)
