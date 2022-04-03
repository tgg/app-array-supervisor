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

type NewModelResponse struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type UpdateResponse struct {
	ComponentId string `json:"componentId"`
	Status      int    `json:"status"`
}

type CommandResponse struct {
	UpdateResponse
	SendCommandRequest
	Result string `json:"result"`
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

func NewNewModelResponse(id string, path string) CustomHubResponse {
	inner := NewModelResponse{id, path}
	return *NewCustomHubResponse(inner, TypeNewModel, StatusOk)
}

func NewUpdateResponse(componentId string, status int) CustomHubResponse {
	return *NewCustomHubResponse(NewUpdateResponseInner(componentId, status), TypeUpdate, StatusOk)
}

func NewUpdateResponseInner(componentId string, status int) UpdateResponse {
	var statusCode int
	if status == 0 {
		statusCode = UpdateOk
	} else {
		statusCode = UpdateNok
	}
	return UpdateResponse{componentId, statusCode}
}

func NewCommandResponse(status int, result string, request SendCommandRequest) CustomHubResponse {
	updateResponse := NewUpdateResponseInner(request.Component, status)
	inner := CommandResponse{
		UpdateResponse:     updateResponse,
		SendCommandRequest: request,
		Result:             result,
	}
	return *NewCustomHubResponse(inner, TypeUpdate, StatusOk)
}
