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

	var content types.Content
	err := json.Unmarshal([]byte(request.Body), &content)

	// Response 500 if Unmarshal failed
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	err = modules.CreateContent(content)

	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	return Response(utils.APIGateway204()), nil
}

func main() {
	lambda.Start(Handler)
}
