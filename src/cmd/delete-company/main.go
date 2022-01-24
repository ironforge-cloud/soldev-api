package main

import (
	"api/internal/database"
	"api/internal/utils"
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Handler AWS Lambda
func Handler(ctx context.Context, request Request) (Response, error) {
	var companyID = request.PathParameters["companyID"]

	db, err := database.GetConnection()
	defer db.Close()

	if err != nil {
		log.Println(err)
		return Response(utils.APIGateway500(errors.New("error connecting to db"))), nil
	}

	err = database.DeleteCompany(db, companyID)
	if err != nil {
		log.Println(err)
		return Response(utils.APIGateway500(errors.New("db error"))), nil
	}

	return Response(utils.APIGateway204()), nil
}

func main() {
	lambda.Start(Handler)
}
