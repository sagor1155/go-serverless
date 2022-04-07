package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var db *dynamodb.DynamoDB

type Todo struct {
	Id      string `json:"id,omitempty"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

func Handler(ctx context.Context, req Request) (Response, error) {
	fmt.Println("Invoked Lambda:", lambdacontext.FunctionName)
	fmt.Println("HTTP Path:", req.Path)
	fmt.Println("HTTP Method:", req.HTTPMethod)

	// get table name from environment variable
	tableName := os.Getenv("DYNAMODB_TABLE")
	fmt.Println("Table Name: ", tableName)

	// create db client
	sess, _ := session.NewSession(&aws.Config{Endpoint: aws.String("http:localhost:8000")})
	db = dynamodb.New(sess)
	fmt.Println("DynamoDB Initialized")

	todos, err := GetTodos(tableName)
	if err != nil {
		return ReturnServerError(http.StatusInternalServerError, "Error saving item into DB", err)
	}

	todosString, err := json.Marshal(todos)
	if err != nil {
		fmt.Println("Got error marshalling result: ", err.Error())
		return ReturnServerError(http.StatusInternalServerError, "Error marshalling result", err)
	}

	return ReturnOK(http.StatusOK, "Successfully saved todo item.", string(todosString))
}

func main() {
	lambda.Start(Handler)
}

func GetTodos(tableName string) ([]Todo, error) {
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	output, err := db.Scan(params)

	if err != nil {
		fmt.Println("Error scanning item from DB: ", err.Error())
		return nil, err
	}

	todos := []Todo{}
	for _, item := range output.Items {
		todo := Todo{}
		err = dynamodbattribute.UnmarshalMap(item, &todo)

		if err != nil {
			fmt.Println("Got error unmarshalling: ", err.Error())
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// type APIGatewayProxyResponse struct {
// 	StatusCode        int                 `json:"statusCode"`
// 	Headers           map[string]string   `json:"headers"`
// 	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
// 	Body              string              `json:"body"`
// 	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
// }

func ReturnOK(status int, message string, data interface{}) (Response, error) {

	body := map[string]interface{}{
		"message": message,
		"status":  "success",
		"error":   "null",
		"data":    data,
	}

	b, _ := json.Marshal(body)

	resp := Response{
		StatusCode: http.StatusOK,
		Body:       string(b),
		Headers:    map[string]string{"accept": "application/json"},
	}

	return resp, nil
}

func ReturnBadRequest(status int, message string, err error) (Response, error) {
	body := map[string]string{
		"message": message,
		"status":  "error",
		"error":   err.Error(),
	}

	b, _ := json.Marshal(body)

	resp := Response{
		StatusCode: http.StatusBadRequest,
		Body:       string(b),
	}

	return resp, nil
}

func ReturnServerError(status int, message string, err error) (Response, error) {
	body := map[string]string{
		"message": message,
		"status":  "error",
		"error":   err.Error(),
	}

	b, _ := json.Marshal(body)

	resp := Response{
		StatusCode: http.StatusInternalServerError,
		Body:       string(b),
	}

	return resp, nil
}
