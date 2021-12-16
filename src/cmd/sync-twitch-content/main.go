package main

import (
	"api/internal/modules"
	"api/internal/utils"
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse

// Handler AWS Lambda
func Handler(ctx context.Context) (Response, error) {
	err := modules.TwitchIntegration()

	// Check if error exist
	if err != nil && err.Error() == "404" {
		return Response(utils.APIGateway404(errors.New("twitch content not found"))), nil
	} else if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	return Response(utils.APIGateway204()), nil
}

func main() {
	lambda.Start(Handler)
}
