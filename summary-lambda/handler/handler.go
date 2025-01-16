package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/supabase-community/supabase-go"
)

type Request struct {
	Id string `json:"docId"`
}

func GetSummary(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	supabaseClient, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "unable to create supabase client"}`),
		}, nil
	}

	var request Request
	err = json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "invalid JSON"}`),
		}, nil
	}

	fmt.Println(request)

	if request.Id == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "docId is required"}`),
		}, nil
	}

	data, _, err  := supabaseClient.From(os.Getenv("TABLE_NAME")).Select("*","exact",false).Eq("id", request.Id).Single().Execute()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "unable to get response from supabase"}`),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(data),
	}, nil
}