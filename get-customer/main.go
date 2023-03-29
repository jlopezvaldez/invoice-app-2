package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type customer struct {
	PhoneNumber   string `json:phone`
	StreetAddress string `json:streetAddress`
	City          string `json:city`
	ZipCode       string `json:zip`
}

var errorLogger = log.New(os.Stderr, "ERROR", log.Llongfile)
var numberRegex = regexp.MustCompile(`[0-9]{10}`)

func showCustomer(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//check if request is phone number
	phoneNumber := req.QueryStringParameters["phoneNumber"]
	if !numberRegex.MatchString(phoneNumber) {
		return clientError(http.StatusBadRequest)
	}

	//fetch a customer from dynamo using phoneNumber
	cus, err := getItem(phoneNumber)
	if err != nil {
		return serverError(err)
	}

	if cus == nil {
		return clientError(http.StatusNotFound)
	}

	js, err := json.Marshal(cus)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

// helper function for errors
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// helper function for errors
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(showCustomer)
}
