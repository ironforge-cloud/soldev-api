package utils

import "time"

// ParseYoutubeTime to Unix timestamp
func ParseTime(youtubeDate string) int64 {

	t, err := time.Parse(time.RFC3339, "2020-03-24T06:44:56Z")
	if err != nil {
		return 0
	}

	return t.Unix()
}
