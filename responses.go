package main

const (
	StatusOk    = 200
	StatusError = 500
)

const (
	TypeError           = "Error"
	TypeExistingModel   = "ExistingModel"
	TypeNewModel        = "NewModel"
	TypeMessage         = "Message"
	TypeUpdate          = "Update"
	TypeCommandResponse = "CommandResponse"
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
	Id    string   `json:"id"`
	Paths []string `json:"paths"`
	Msg   string   `json:"msg"`
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

func NewNewModelResponse(id string, paths []string) CustomHubResponse {
	inner := NewModelResponse{id, paths, ""}
	return *NewCustomHubResponse(inner, TypeNewModel, StatusOk)
}

func NewExistingModelResponse(id string, paths []string, msg string) CustomHubResponse {
	inner := NewModelResponse{id, paths, msg}
	return *NewCustomHubResponse(inner, TypeExistingModel, StatusOk)
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
	return *NewCustomHubResponse(inner, TypeCommandResponse, StatusOk)
}
