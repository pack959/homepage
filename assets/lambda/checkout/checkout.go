package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"checkout/stripe"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	AWSLambdaFunctionVersion = "AWS_LAMBDA_FUNCTION_VERSION"
	StripeApiKey             = "STRIPE_KEY_SECRET"
)

type data struct {
	PublicURL   string `json:"public_url"`
	SuccessPath string `json:"success_path"`
	CancelPath  string `json:"cancel_path"`
	Items       []struct {
		Quantity    int    `json:"quantity"`
		Description string `json"description"`
		Amount      int    `json:"amount"`
	}
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 501,
		}, nil
	}

	data := new(data)
	if err := json.Unmarshal([]byte(request.Body), data); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil
	}

	// Get secret key
	key, found := os.LookupEnv(StripeApiKey)
	if !found {
		return nil, errors.New("environment variable not set")
	}
	processor := stripe.New(data.PublicURL, key)

	session, err := processor.CreateCheckoutSession(data.SuccessPath, data.CancelPath)
	if err != nil {
		return nil, err
	}
	session.AddItem("abc", "def", 123, 1)
	session.Start()

	// response := map[string]interface{}{
	// 	"session_id": session.GetID(),
	// }
	// responseJSON, _ := json.Marshal(response)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       session.GetID(),
	}, nil
}

func main() {
	_, ok := os.LookupEnv(AWSLambdaFunctionVersion)
	if ok {
		log.Printf("Running in AWS lambda environment, starting lambda handler.")
		lambda.Start(handler)
		os.Exit(0)
	}

	log.Printf("Not running in AWS lambda environment, starting mock handler.")
	os.Exit(-1)
}
