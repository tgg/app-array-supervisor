package main

import "encoding/json"

const (
	StatusOk    = 200
	StatusError = 500
)

const (
	TypeError    = "Error"
	TypeNewModel = "NewModel"
)

type CustomHubResponse struct {
	msg        string
	jsonType   string
	statusCode int
}

func NewCustomHubResponse(message string, jsonType string, statusCode int) *CustomHubResponse {
	return &CustomHubResponse{
		statusCode: statusCode,
		jsonType:   jsonType,
		msg:        message,
	}
}

func NewErrorResponse(message string) CustomHubResponse {
	return *NewCustomHubResponse(message, TypeError, StatusError)
}

type NewModelResponse struct {
	id   string
	path []string
}

func NewNewModelResponse(id string, path string) CustomHubResponse {
	inner := NewModelResponse{}
	b, _ := json.Marshal(inner)
	return *NewCustomHubResponse(string(b), TypeNewModel, StatusOk)
}
