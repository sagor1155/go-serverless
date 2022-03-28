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

	// parse request
	todo := Todo{}
	err := json.Unmarshal([]byte(req.Body), &todo)
	if err != nil {
		fmt.Println("Invalid data format!!")
		return ReturnBadRequest(http.StatusBadRequest, "Invalid data format!!", err)
	}

	fmt.Println("ID:", todo.Id)
	fmt.Println("Title:", todo.Title)
	fmt.Println("Details:", todo.Details)

	// get table name from environment variable
	tableName := os.Getenv("DYNAMODB_TABLE")
	fmt.Println("Table Name: ", tableName)

	// create db client
	sess, _ := session.NewSession()
	db = dynamodb.New(sess)
	fmt.Println("DynamoDB Initialized")

	// update dynamo db
	err = CreateTodo(todo, tableName)
	if err != nil {
		return ReturnServerError(http.StatusInternalServerError, "Error saving item into DB", err)
	}

	return ReturnOK(http.StatusOK, "Successfully saved todo item.", todo)
}

func main() {
	lambda.Start(Handler)
}

func CreateTodo(todo Todo, tableName string) error {
	// update dynamo db
	inputAttr, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		fmt.Println("Error marshalling item: ", err.Error())
		// return ReturnServerError(http.StatusInternalServerError, "Error marshalling item", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      inputAttr,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)
	if err != nil {
		fmt.Println("Error saving item into DB: ", err.Error())
		// return ReturnServerError(http.StatusInternalServerError, "Error saving item into DB", err)
		return err
	}

}

func ReturnOK(status int, message string, data interface{}) (Response, error) {
	// generate response
	body := map[string]string{
		"message": message,
		"status":  "success",
		"error":   "null",
	}

	b, _ := json.Marshal(body)

	resp := Response{
		StatusCode: http.StatusOK,
		Body:       string(b),
	}

	return resp, nil
}

func ReturnBadRequest(status int, message string, err error) (Response, error) {
	// generate response
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
	// generate response
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
