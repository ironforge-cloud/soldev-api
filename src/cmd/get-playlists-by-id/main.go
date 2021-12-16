package main

import (
	"api/internal/modules"
	"api/internal/utils"
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse

// Request represents an API Gateway request
type Request events.APIGatewayProxyRequest

func Handler(ctx context.Context, request Request) (Response, error) {
	var vertical = request.PathParameters["vertical"]
	var playlistID = request.PathParameters["playlistID"]

	playlist, err := modules.GetPlaylistByID(vertical, playlistID)

	// Check if error exist
	if err != nil && err.Error() == "404" {
		return Response(utils.APIGateway404(errors.New("playlist not found"))), nil
	} else if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	// Marshal the response into json bytes
	response, err := json.Marshal(&playlist)
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	return Response(utils.APIGateway200(response)), nil
}

func main() {
	lambda.Start(Handler)
}
