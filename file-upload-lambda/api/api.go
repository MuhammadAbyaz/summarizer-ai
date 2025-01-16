package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"lambda/s3"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type APIHandler struct {
	s3Client s3.S3Client
}

func NewAPIHandler(client s3.S3Client) APIHandler {
	return APIHandler{
		s3Client: client,
	}
}

type RequestBody struct {
	File string `json:"file"` 
}

func (u APIHandler) FileUploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body RequestBody
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Invalid JSON: %v", err),
		}, nil
	}
	if body.File == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "File content is required",
		}, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(body.File)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Failed to decode base64: %v", err),
		}, nil
	}

	filename := fmt.Sprintf("upload_%d", time.Now().Unix())
	
	_, err = u.s3Client.UploadFile(filename, decoded)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("successfully uploaded %s", filename),
	}, nil
}
