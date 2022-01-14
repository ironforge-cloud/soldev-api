package types

// Content ...
type Content struct {
	// DynamoDB composite key
	PK string
	SK string

	ContentStatus string // "active", "inactive" and "submitted"

	Position   int64 // content != playlists means weight
	Tags       []string
	Lists      string // TODO: in the future this will be a slice
	SpecialTag string // "new", "hot", "best" and "old"

	Url         string
	Title       string
	Description string
	PublishedAt string
	Img         string
	Author      string
	ContentType string
	Vertical    string

	// Video specific
	PlaylistID string
	Promoted   int8
	Live       int8
	Provider   string
	Expdate    int64

	PlaylistTitle string
}

// Twitch ...
type Twitch struct {
	ID          string `json:"id"`
	Channel     string `json:"user_name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Img         string `json:"thumbnail_url"`
	StartedAt   string `json:"started_at"`
	CreatedAt   string `json:"created_at"`
	Url         string `json:"url"`
}

// TwitchResponse ...
type TwitchResponse struct {
	Data []Twitch `json:"data"`
}

type AlgoliaRecord struct {
	ObjectID string `json:"objectID"`
	Content
}
