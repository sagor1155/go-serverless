package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

func Handler(ctx context.Context, req Request) (Response, error) {
	fmt.Println("Invoked Lambda:", lambdacontext.FunctionName)
	fmt.Println("HTTP Path:", req.Path)
	fmt.Println("HTTP Method:", req.HTTPMethod)

	body := map[string]string{
		"message": "Lambda executed successfully",
		"status":  "success",
	}

	b, _ := json.Marshal(body)

	resp := Response{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 200,
		Body:       string(b),
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
