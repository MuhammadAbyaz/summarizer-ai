package main

import (
	"lambda/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main(){
	lambda.Start(handler.FileUploadHandler)
}