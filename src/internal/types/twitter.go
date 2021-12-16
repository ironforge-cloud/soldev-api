package types

type TwitterResponse struct {
	Data     []Tweet  `json:"data"`
	Includes Includes `json:"includes"`
	Meta     Meta     `json:"meta"`
}

type Tweet struct {
	Text            string            `json:"text"`
	CreatedAt       string            `json:"created_at"`
	ID              string            `json:"id,omitempty"`
	PublicMetrics   PublicMetrics     `json:"public_metrics"`
	Attachments     Attachments       `json:"attachments"`
	AuthorID        string            `json:"author_id"`
	ReferencedTweet []ReferencedTweet `json:"referenced_tweets"`
}

type PublicMetrics struct {
	RetweetCount int `json:"retweet_count"`
	ReplyCount   int `json:"reply_count"`
	LikeCount    int `json:"like_count"`
	QuoteCount   int `json:"quote_count"`
}

type Includes struct {
	Users  []TwitterUser `json:"users"`
	Tweets []Tweet       `json:"tweets"`
	Media  []Media       `json:"media"`
}

type TwitterUser struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	ProfileImagineUrl string `json:"profile_image_url"`
	Name              string `json:"name"`
	Url               string `json:"url"`
	Verified          bool   `json:"verified"`
	Protected         bool   `json:"protected"`
}

type Attachments struct {
	MediaKeys []string `json:"media_keys"`
}

type ReferencedTweet struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Media struct {
	Type            string `json:"type,omitempty"`
	PreviewImageUrl string `json:"preview_image_url"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	MediaKey        string `json:"media_key"`
}

type Meta struct {
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

type TwitterListResponseData struct {
	PK               string
	Media            []Media
	ReferencedTweet  []ReferencedTweet
	ReferencedTweets []TwitterListResponseData
	Author           TwitterUser
	Text             string        `json:"text"`
	CreatedAt        string        `json:"created_at"`
	ID               string        `json:"id,omitempty"`
	PublicMetrics    PublicMetrics `json:"public_metrics"`
	Attachments      Attachments   `json:"attachments"`
	AuthorID         string        `json:"author_id"`
	Type             string
	Expdate          int64
	Pinned           int8
}
