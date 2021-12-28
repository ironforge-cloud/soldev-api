package utils

import (
	"fmt"
	"hash/crc32"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

// APIGateway500 responds status code 500 using APIGatewayProxyResponse
func APIGateway500(err error) Response {
	return Response{
		Body:       err.Error(),
		StatusCode: 500,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Max-Age:":          "604800",
		},
	}
}

// APIGateway404 responds status code 404 Not Found using APIGatewayProxyResponse
func APIGateway404(err error) Response {
	return Response{
		Body:       err.Error(),
		StatusCode: 404,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Max-Age:":          "604800",
		},
	}
}

// APIGateway204 responds status code 204 Not Content using APIGatewayProxyResponse
func APIGateway204() Response {
	return Response{
		StatusCode: 204,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Max-Age:":          "604800",
		},
	}
}

// APIGateway200 responds status code 200 and a payload using APIGatewayProxyResponse
func APIGateway200(data []byte) Response {
	// Helper for UI empty responses
	body := ""
	if len(data) == 0 {
		body = "[]"
	} else {
		body = string(data)
	}

	return Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Etag":                             etag(data),
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Max-Age:":          "604800",
			"Cache-Control":                    "public, max-age=1, s-maxage=1, proxy-revalidate, must-revalidate",
		},
		Body: body,
	}
}

// etag generator
func etag(data []byte) string {
	crc := crc32.ChecksumIEEE(data)
	return fmt.Sprintf(`W/"%s-%d-%08X"`, len(data), crc)
}
