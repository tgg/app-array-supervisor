package main

const (
	StatusOk    = 200
	StatusError = 500
)

const (
	TypeError    = "Error"
	TypeNewModel = "NewModel"
)

type CustomHubResponse struct {
	Msg        any    `json:"msg"`
	JsonType   string `json:"type"`
	StatusCode int    `json:"statusCode"`
}

func NewCustomHubResponse(message any, jsonType string, statusCode int) *CustomHubResponse {
	return &CustomHubResponse{
		StatusCode: statusCode,
		JsonType:   jsonType,
		Msg:        message,
	}
}

func NewErrorResponse(message string) CustomHubResponse {
	return *NewCustomHubResponse(message, TypeError, StatusError)
}

type NewModelResponse struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

func NewNewModelResponse(id string, path string) CustomHubResponse {
	inner := NewModelResponse{id, path}
	return *NewCustomHubResponse(inner, TypeNewModel, StatusOk)
}
