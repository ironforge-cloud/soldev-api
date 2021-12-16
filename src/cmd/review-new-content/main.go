package main

import (
	"api/internal/database"
	"api/internal/modules"
	"api/internal/utils"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse

// Request represents an API Gateway request
type Request events.APIGatewayProxyRequest

func Handler(ctx context.Context, request Request) (Response, error) {

	// Get list of content with new specialTag
	content, err := database.QueryContent("", "", "New")

	if err != nil {
		log.Printf("Error reviewing new content: %v ", err)
		return Response(utils.APIGateway500(err)), nil
	}

	// Review content
	contentList, err := modules.ReviewNewContent(content)

	if err != nil {
		log.Printf("Error reviewing new content: %v ", err)
		return Response(utils.APIGateway500(err)), nil
	}

	// Save content
	for _, item := range contentList {
		err := database.SaveContent(item)

		if err != nil {
			log.Printf("Error while updating reviewed content: %v ", err)
			return Response(utils.APIGateway500(err)), nil
		}
	}

	return Response(utils.APIGateway204()), nil
}

func main() {
	lambda.Start(Handler)
}
