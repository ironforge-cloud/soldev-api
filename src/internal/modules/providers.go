package modules

import (
	"api/internal/database"
	"api/internal/providers"
	"api/internal/types"
	"errors"
)

// YoutubeIntegration find all YouTube playlists in our DB, using the
// YouTube API fetch all the information and then saves the information
// in the DB
func YoutubeIntegration() error {
	// Get YouTube Playlists
	playlists, err := database.GetPlaylistsByProvider("Youtube")
	if err != nil {
		return err
	}

	var content []types.Content
	// Find videos for each playlists
	for _, playlist := range playlists {

		youtubeData, _ := providers.GetYoutubePlaylist(playlist)

		content = append(content, youtubeData...)
	}

	// If there are no content, return an error
	if content == nil {
		return errors.New("404")
	}

	// Save the videos of each playlist in the content table
	for _, item := range content {
		err = database.SaveContent(item)
		if err != nil {
			return err
		}
	}

	return nil
}
