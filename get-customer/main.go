package main

import "github.com/aws/aws-lambda-go/lambda"

type customer struct {
	PhoneNumber   string `json:phone`
	StreetAddress string `json:streetAddress`
	City          string `json:city`
	ZipCode       string `json:zip`
}

func showCustomer() (*customer, error) {
	//fetch a customer from dynamo using phoneNumber
	cus, err := getItem("8888888888")
	if err != nil {
		return nil, err
	}

	return cus, nil
}

func main() {
	lambda.Start(showCustomer)
}
