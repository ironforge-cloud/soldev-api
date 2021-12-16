package main

import (
	"api/internal/modules"
	"api/internal/utils"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse

// Request represents an API Gateway request
type Request events.APIGatewayProxyRequest

func Handler(ctx context.Context, request Request) (Response, error) {
	var vertical = request.PathParameters["vertical"]
	var contentType = request.PathParameters["contentType"]
	var tags = request.QueryStringParameters["tags"]
	var specialTags = request.QueryStringParameters["specialTags"]

	content, err := modules.GetContent(vertical, contentType, tags, specialTags)

	// Check if error exist
	if err != nil && err.Error() == "404" {
		return Response(utils.APIGateway200([]byte{})), nil
	} else if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	// Marshal the response into json bytes
	response, err := json.Marshal(&content)
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	return Response(utils.APIGateway200(response)), nil
}

func main() {
	lambda.Start(Handler)
}
