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
	webHookURL = "MANUAL_DEPLOY_WEBHOOK"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Get secret key
	url, found := os.LookupEnv(webHookURL)
	if !found {
		return nil, errors.New("environment variable not set")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return nil, errors.New("error posting to webhook")
	}
	defer resp.Body.Close()

	headers := map[string]string{
		"Location": "/_internal/site/success/",
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
