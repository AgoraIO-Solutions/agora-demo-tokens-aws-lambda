package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleLambdaEvent() (events.APIGatewayProxyResponse, error) {
	bodyText := fmt.Sprintf("Hello World ")

	return events.APIGatewayProxyResponse{
		Body:       bodyText,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
