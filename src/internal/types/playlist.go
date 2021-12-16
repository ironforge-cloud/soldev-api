package types

// Playlist represents a content playlist
type Playlist struct {
	ID          string
	Provider    string
	Title       string
	Description string
	Author      string
	Vertical    string
	Tags        []string
	Position    int
}
