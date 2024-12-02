package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ReleaseDateLayout = "02.01.2006"

type ReleaseDate time.Time

func (t *ReleaseDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("InfoResponse.UnmarshalJSON: input is not a JSON string")
	}

	value := strings.Trim(string(data), `"`)
	date, err := time.Parse(ReleaseDateLayout, value)
	if err != nil {
		return fmt.Errorf("InfoResponse.UnmarshalJSON: %w ", err)
	}

	*t = ReleaseDate(date)
	return nil
}

func (t *ReleaseDate) String() string {
	return time.Time(*t).String()
}
