package main

const (
	StatusOk    = 200
	StatusError = 500
)

const (
	TypeError    = "Error"
	TypeNewModel = "NewModel"
	TypeMessage  = "Message"
	TypeUpdate   = "Update"
)

const (
	UpdateOk  = 0
	UpdateNok = 1
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

func NewMessageResponse(message string) CustomHubResponse {
	return *NewCustomHubResponse(message, TypeMessage, StatusOk)
}

type NewModelResponse struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

func NewNewModelResponse(id string, path string) CustomHubResponse {
	inner := NewModelResponse{id, path}
	return *NewCustomHubResponse(inner, TypeNewModel, StatusOk)
}

type UpdateResponse struct {
	ComponentId string `json:"componentId"`
	Status      int    `json:"status"`
}

func NewUpdateResponse(componentId string, status int) CustomHubResponse {
	var statusCode int
	if status == 0 {
		statusCode = UpdateOk
	} else {
		statusCode = UpdateNok
	}
	inner := UpdateResponse{componentId, statusCode}
	return *NewCustomHubResponse(inner, TypeUpdate, StatusOk)
}
