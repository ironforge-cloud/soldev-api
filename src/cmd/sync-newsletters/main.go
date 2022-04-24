package main

import (
	"api/internal/modules"
	"api/internal/types"
	"api/internal/utils"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Handler AWS Lambdad
func Handler(ctx context.Context, request Request) (Response, error) {

	var data []types.Content

	count := 0

	for i := 2; i >= 0; i-- {
		response := utils.FindData(i)

		buf := len(response)
		for _, item := range response {
			data = append(data, types.Content{
				ContentType:     "newsletters",
				Vertical:        "Solana",
				PublishedAt:     item.DateAdded,
				ContentMarkdown: item.ContentMarkdown,
				Title:           item.Title,
				Description:     item.Brief,
				SK:              item.Slug,
				PK:              "Solana#newsletters",
				Img:             item.Img,
				ContentStatus:   "active",
				Url:             "https://soldev.app/newsletters/" + item.Slug,
				Position:        int64(count + buf),
			})

			buf--
		}

		count += len(response)
	}

	err := modules.SaveContent(data)

	// Respond 500 if error occurred while saving the playlists in the database
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	return Response(utils.APIGateway204()), nil
}

func main() {
	lambda.Start(Handler)
}
