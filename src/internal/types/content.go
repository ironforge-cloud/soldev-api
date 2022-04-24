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

	// Newsletter specific
	ContentMarkdown string

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

type HashNodeUser struct {
	User map[string]Publication
}

type Publication struct {
	Posts []HashnodePost `json:"posts"`
}

type HashnodePost struct {
	Title           string `json:"title"`
	Brief           string `json:"brief"`
	Slug            string `json:"slug"`
	DateAdded       string `json:"dateAdded"`
	ContentMarkdown string `json:"contentMarkdown"`
	Img             string `json:"coverImage"`
}

// map[user:map[publication:map[posts:[map
// map[user:{map[]}]
