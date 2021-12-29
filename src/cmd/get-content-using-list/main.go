package main

import (
	"api/internal/modules"
	"api/internal/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse

// Request represents an API Gateway request
type Request events.APIGatewayProxyRequest

func Handler(ctx context.Context, request Request) (Response, error) {
	var list = request.PathParameters["listName"]

	// Get list of content with specialTag
	queryResponse, err := modules.FilterContentByList(list)

	// Check if error exist
	if err != nil && err.Error() == "404" {
		return Response(utils.APIGateway200([]byte{})), nil
	} else if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	// Marshal the response into json bytes
	response, err := json.Marshal(&queryResponse)
	fmt.Println("Response", response)
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	return Response(utils.APIGateway200(response)), nil
}

func main() {
	lambda.Start(Handler)
}
