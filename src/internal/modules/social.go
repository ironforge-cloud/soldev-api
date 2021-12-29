package modules

import (
	"api/internal/database"
	"api/internal/providers"
	"api/internal/types"
	"time"
)

// SyncTwitterList ..
func SyncTwitterList() error {

	tweets, err := providers.TwitterList("1444990678371651585")
	if err != nil {
		return err
	}

	err = database.SaveTweets(tweets, false)
	if err != nil {
		return err
	}

	return nil
}

// GetTweets finds all the Tweets for the Twitter List 1444990678371651585
func GetTweets() ([]types.TwitterListResponseData, error) {
	tweets, err := database.QueryTweets("1444990678371651585")
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

// PinTweet find tweets using tweet id
func PinTweet(id string) error {
	// Find Tweet
	tweets, err := database.QueryTweetsUsingTweetID(id)
	if err != nil {
		return err
	}

	// Update Data
	for i, tweet := range tweets {
		if tweet.Pinned == 0 {
			tweets[i].Pinned = 1
			// Removing TTL from DynamoDB
			tweets[i].Expdate = 0
		} else {
			tweets[i].Pinned = 0
			// Adding TTL for DynamoDB
			tweets[i].Expdate = time.Now().Add(time.Hour * 24).Unix()
		}
	}

	// Save Tweet
	err = database.SaveTweets(tweets, true)
	if err != nil {
		return err
	}

	return nil
}
