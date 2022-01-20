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
	Promoted   int8 // deprecated
	Live       int8 // deprecated
	Provider   string
	Expdate    int64

	PlaylistTitle string
}

type AlgoliaRecord struct {
	ObjectID string `json:"objectID"`
	Content
}
