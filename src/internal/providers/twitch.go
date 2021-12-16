package providers

import (
	"api/internal/types"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

var accessToken string

// TwitchClient represents a new TwitchClient
func TwitchClient() string {
	oauth2Config := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return token.AccessToken
}

// TwitchRequest ...
func TwitchRequest(endpoint string, typeOfID string, ID string) ([]byte, error) {
	if accessToken == "" {
		accessToken = TwitchClient()
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("TWITCH_HELIX_URL")+endpoint, nil)
	if err != nil {
		log.Printf("Error while requesting streams to Twitch: %v", err)
		return nil, err
	}

	// Headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))

	// Query Params
	query := req.URL.Query()
	query.Add(typeOfID, ID)
	query.Add("first", "100")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error while doing client.Do(): %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// If response is Unauthorized we need to create a new Bearer token
	// because the previous expired
	if resp.StatusCode == 401 {
		accessToken = TwitchClient()
	}

	// Convert response to bytes
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return bodyBytes, nil
}

// FetchSavedStreams fetches all saved streams
func FetchSavedStreams(userId string) ([]types.Content, error) {
	bodyBytes, err := TwitchRequest("videos", "user_id", userId)
	if err != nil {
		return nil, err
	}

	// Structure to convert bytes
	var twitchResponse types.TwitchResponse
	// Unmarshal bytes into a iterable structure
	err = json.Unmarshal(bodyBytes, &twitchResponse)
	if err != nil {
		return nil, err
	}

	var content []types.Content
	for _, item := range twitchResponse.Data {

		// Thumbnail size needs to be added it to the URL
		imgWithSize := strings.Replace(item.Img, "%{width}", "360", -1)
		imgWithSize = strings.Replace(imgWithSize, "%{height}", "202", -1)

		video := types.Content{
			PK:            "Solana#twitch-solana",
			SK:            item.ID,
			ContentStatus: "active",
			Url:           item.Url,
			Title:         item.Title,
			PublishedAt:   item.CreatedAt,
			Img:           imgWithSize,
			Author:        item.Channel,
			ContentType:   "Playlist",
			Vertical:      "Solana",
			PlaylistID:    "twitch-solana",
			Promoted:      0,
			Live:          0,
			Provider:      "Twitch",
			Expdate:       time.Now().Add(time.Hour * 12).Unix(),
		}

		content = append(content, video)
	}

	return content, nil
}

// FetchLiveStream ...
func FetchLiveStream(userID string) ([]types.Content, error) {
	bodyBytes, err := TwitchRequest("streams", "user_id", userID)
	if err != nil {
		return nil, err
	}

	// Structure to convert bytes
	var twitchResponse types.TwitchResponse
	// Unmarshal bytes into an iterable structure
	err = json.Unmarshal(bodyBytes, &twitchResponse)
	if err != nil {
		return nil, err
	}

	// If the stream did not start, we redirect users to twitch channel
	if len(twitchResponse.Data) == 0 {
		twitchResponse.Data = append(twitchResponse.Data, types.Twitch{
			Title: "Livestream",
			Img:   "/solana-logo.jpg",
		})
	}

	var content []types.Content
	for _, item := range twitchResponse.Data {
		// Thumbnail size needs to be added it to the URL
		imgWithSize := strings.Replace(item.Img, "{width}", "360", -1)
		imgWithSize = strings.Replace(imgWithSize, "{height}", "202", -1)

		// TODO: Quick hack for the UI, please fix this shit
		var live int8 = 1
		if item.Img == "/solana-logo.jpg" {
			live = 0
		}

		video := types.Content{
			PK:            "Solana#twitch-solana",
			SK:            "livestream",
			ContentStatus: "active",
			Url:           "https://www.twitch.tv/solanatv", // TODO: hardcoded for now
			Title:         item.Title,
			PublishedAt:   item.StartedAt,
			Img:           imgWithSize,
			Author:        "SolanaTV",
			ContentType:   "Livestream",
			Vertical:      "Solana",
			PlaylistID:    "twitch-solana",
			Promoted:      1,
			Live:          live,
			Provider:      "Twitch",
			Expdate:       time.Now().Add(time.Minute).Unix(),
		}

		content = append(content, video)
	}

	return content, nil
}
