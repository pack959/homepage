package main

import (
	"bytes"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	WebhookURL = "MANUAL_DEPLOY_WEBHOOK"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Get secret key
	url, found := os.LookupEnv(WebhookURL)
	if !found {
		return nil, errors.New("environment variable not set")
	}

	// ignore all errors and responses
	http.Post(url, "application/json", bytes.NewBuffer([]byte(`{}`)))

	headers := map[string]string{
		"Location": "/",
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
