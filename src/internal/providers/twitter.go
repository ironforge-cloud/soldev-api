package providers

import (
	"api/internal/types"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func TwitterList(listID string) ([]types.TwitterListResponseData, error) {
	url := "https://api.twitter.com/2/lists/" + listID + "/tweets"

	//  Twitter Auth Token
	var bearer = "Bearer " + os.Getenv("TWITTER_BEARER_TOKEN")

	// Create a new req
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// add params
	q := req.URL.Query()
	q.Add("expansions", "author_id,attachments.media_keys,referenced_tweets.id,referenced_tweets.id.author_id")
	q.Add("tweet.fields", "attachments,author_id,public_metrics,created_at,id,in_reply_to_user_id,referenced_tweets,text")
	q.Add("user.fields", "id,name,profile_image_url,protected,url,username,verified")
	q.Add("media.fields", "duration_ms,height,media_key,preview_image_url,type,url,width,public_metrics")
	q.Add("max_results", "10")
	req.URL.RawQuery = q.Encode()

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var response types.TwitterResponse

	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	//  Format data to be saved in DB
	formattedData := formatTwitterResponse(response, listID)

	return formattedData, nil
}

func formatTwitterResponse(response types.TwitterResponse, listID string) []types.TwitterListResponseData {
	var data []types.TwitterListResponseData
	for _, tweet := range response.Data {
		var author types.TwitterUser
		var fullReferencedTweets []types.TwitterListResponseData
		var media []types.Media

		// Find Tweet Author data
		author = findAuthor(tweet.AuthorID, response.Includes.Users)

		// Find Referenced Tweets data if any
		if tweet.ReferencedTweet != nil {
			fullReferencedTweets = findReferencedTweets(tweet.ReferencedTweet, response.Includes.Tweets, response.Includes.Users, listID)
		}

		// Find Tweet Media if any
		if tweet.Attachments.MediaKeys != nil {
			media = findMedia(tweet.Attachments, response.Includes.Media)
		}

		// Define Final Data Structure
		data = append(data, types.TwitterListResponseData{
			PK:               listID,
			Media:            media,
			Author:           author,
			ReferencedTweets: fullReferencedTweets,
			Text:             tweet.Text,
			CreatedAt:        tweet.CreatedAt,
			ID:               tweet.ID,
			PublicMetrics:    tweet.PublicMetrics,
			Attachments:      tweet.Attachments, // this may not be needed because Media
			AuthorID:         tweet.AuthorID,    // this may not be needed because Author
			Expdate:          time.Now().Add(time.Hour * 24).Unix(),
		})
	}

	return data
}

// Finds Tweet Author data using AuthorID
func findAuthor(authorID string, data []types.TwitterUser) types.TwitterUser {
	var response types.TwitterUser
	for _, user := range data {
		if authorID != user.ID {
			continue
		}

		response = user
		break
	}

	return response
}

func findReferencedTweets(referencedTweets []types.ReferencedTweet, tweets []types.Tweet, users []types.TwitterUser, listID string) []types.TwitterListResponseData {
	var fullReferencedTweets []types.TwitterListResponseData
	for _, referencedTweet := range referencedTweets {
		for _, tweet := range tweets {
			if referencedTweet.ID != tweet.ID {
				continue
			}

			data := types.TwitterListResponseData{
				PK:              listID,
				Text:            tweet.Text,
				CreatedAt:       tweet.CreatedAt,
				ID:              tweet.ID,
				PublicMetrics:   tweet.PublicMetrics,
				Attachments:     tweet.Attachments,
				AuthorID:        tweet.AuthorID,
				ReferencedTweet: tweet.ReferencedTweet,
				Author:          findAuthor(tweet.AuthorID, users),
				Type:            referencedTweet.Type,
				Expdate:         time.Now().Add(time.Hour * 24).Unix(),
			}

			fullReferencedTweets = append(fullReferencedTweets, data)
		}
	}

	return fullReferencedTweets
}

func findMedia(attachments types.Attachments, mediaList []types.Media) []types.Media {
	var media []types.Media
	for _, mediaKey := range attachments.MediaKeys {
		for _, mediaItem := range mediaList {
			if mediaKey == mediaItem.MediaKey {
				media = append(media, mediaItem)
			}
		}
	}

	return media
}
