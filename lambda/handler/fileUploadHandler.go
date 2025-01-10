package handler

import (
	"fmt"
	"lambda/types"
)


func FileUploadHandler(event types.Event) (string, error) {
	if event.Username == ""{
		return "", fmt.Errorf("username cannot be empty")
		
	}
	return fmt.Sprintf("Lambda called by %s", event.Username), nil
}