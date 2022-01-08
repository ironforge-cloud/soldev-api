package providers

import (
	"api/internal/types"
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YoutubeClient represents a new YouTube Service
func YoutubeClient() *youtube.Service {
	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("Error creating new Youtube Client: %v", err)
	}

	return youtubeService
}

// GetYoutubePlaylist finds data for each playlist
func GetYoutubePlaylist(playlist types.Playlist) ([]types.Content, error) {
	client := YoutubeClient()

	part := []string{
		"snippet",
	}

	call := client.PlaylistItems.List(part)
	call = call.MaxResults(100) // 100 videos in each playlist max
	call = call.PlaylistId(playlist.ID)
	response, err := call.Do()

	if err != nil {
		log.Printf("Error while getting youtube playlist %v: %v", playlist.ID, err)
		return nil, err
	}

	var contentList []types.Content
	for _, item := range response.Items {

		// YouTube API can return some private or deleted videos that we need to ignore
		if item.Snippet.Title == "Private video" || item.Snippet.Title == "Deleted video" {
			continue
		}

		contentList = append(contentList, types.Content{
			PK:            playlist.Vertical + "#" + playlist.ID,
			SK:            item.Snippet.ResourceId.VideoId,
			ContentStatus: "active",
			Url:           "https://www.youtube.com/watch?v=" + item.Snippet.ResourceId.VideoId,
			Title:         item.Snippet.Title,
			Description:   item.Snippet.Description,
			PublishedAt:   item.Snippet.PublishedAt,
			Img:           item.Snippet.Thumbnails.High.Url,
			Author:        item.Snippet.ChannelTitle,
			Position:      item.Snippet.Position,
			ContentType:   "Playlist",
			Vertical:      playlist.Vertical,
			PlaylistID:    playlist.ID,
			Promoted:      checkPromote(item.Snippet.ResourceId.VideoId),
			PlaylistTitle: playlist.Title,
			Live:          0,
			Provider:      playlist.Provider,
			// after 24 hours if this video did not come in another api call, we should deleted it
			Expdate: time.Now().Add(time.Hour * 24).Unix(),
		})
	}

	return contentList, nil
}

func checkPromote(contentID string) int8 {
	promotedList := []string{"itWOUcsMjT4", "scorFcEXZD8", "cvW8EwGHw8U", "QcXcqIYflv8", "atAYmR32eeE", "sOlq1S6e1HE", "r7D2CA8qPHY", "vbESNVOZflU"}

	for _, ID := range promotedList {
		if ID == contentID {
			return 1
		}
	}

	return 0
}
