package modules

import (
	"api/internal/database"
	"api/internal/types"
)

// GetPlaylists return all playlists saved
func GetPlaylists(vertical string) ([]types.Playlist, error) {
	playlists, err := database.GetAllPlaylists(vertical)

	if err != nil {
		return nil, err
	}

	return playlists, nil
}

// SavePlaylists saves one or multiple playlists
func SavePlaylists(playlists []types.Playlist) error {
	err := database.SavePlaylists(playlists)

	if err != nil {
		return err
	}

	return nil
}

// GetPlaylistByID ...
func GetPlaylistByID(vertical string, ID string) (types.Playlist, error) {
	playlist, err := database.GetPlaylistByID(vertical, ID)

	if err != nil {
		return types.Playlist{}, err
	}

	return playlist, nil
}
