package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"lambda/s3"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
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
			Body:       fmt.Sprintf("{'error': 'invalid JSON'}"),
		}, nil
	}
	if body.File == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("{'error': 'file content is required'}"),
		}, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(body.File)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("{'error': 'Failed to decode base64'}"),
		}, nil
	}

	filename := uuid.New().String()

	_, err = u.s3Client.UploadFile(filename, decoded)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("{'error': 'cannot upload file to s3'}"),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"id": "%v"}`, filename),
	}, nil
}
