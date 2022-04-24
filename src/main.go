package main

import (
	"api/internal/types"
	"api/internal/utils"
	"fmt"
)

func main() {
	var data []types.Content

	count := 0
	for i := 2; i >= 0; i-- {
		response := utils.FindData(i)

		for _, item := range response {
			data = append(data, types.Content{
				ContentType:     "newsletters",
				Vertical:        "solana",
				PublishedAt:     item.DateAdded,
				ContentMarkdown: item.ContentMarkdown,
				Title:           item.Title,
				Description:     item.Brief,
				SK:              item.Slug,
				PK:              "solana#newsletters",
				Img:             item.Img,
				ContentStatus:   "active",
				Position:        int64(count),
			})

			count++
		}

		fmt.Println(data[0])
	}
}
