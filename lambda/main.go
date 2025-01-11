package main

import (
	"lambda/api"
	"lambda/s3"

	"github.com/aws/aws-lambda-go/lambda"
)

func main(){
	s3Client := s3.NewS3Client()
	apiHandler := api.NewAPIHandler(s3Client)
	lambda.Start(apiHandler.FileUploadHandler)
}