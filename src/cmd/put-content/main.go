package main

import (
	"api/internal/modules"
	"api/internal/types"
	"api/internal/utils"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Handler AWS Lambda
func Handler(ctx context.Context, request Request) (Response, error) {

	var contentPlaylistData []types.Content
	err := json.Unmarshal([]byte(request.Body), &contentPlaylistData)

	// Response 500 if Unmarshal failed
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	err = modules.SaveContent(contentPlaylistData)

	// Respond 500 if error occurred while saving the playlists in the database
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	return Response(utils.APIGateway204()), nil
}

func main() {
	lambda.Start(Handler)
}
