package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/supabase-community/supabase-go"
)

type Request struct {
	Id string `json:"id"`
}

func GetSummary(event events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error){
	supabaseClient,err := supabase.NewClient(os.Getenv("SUPABASE_URL"),os.Getenv("SUPABASE_KEY"),nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body: "unable to create supabase client",
		}, err
	}
	var request Request
	err = json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body: "unable to parse the body",
		}, err
	}
	data, _, err := supabaseClient.From(os.Getenv("TABLE_NAME")).Select("*","exact",false).Single().Eq("id",request.Id).Execute()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body: "unable to get response from supabase",
		}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(data)},nil
}