package main

const (
	StatusOk    = 200
	StatusError = 500
)

const (
	TypeError                   = "Error"
	TypeExistingModel           = "ExistingModel"
	TypeNewModel                = "NewModel"
	TypeMessage                 = "Message"
	TypeInfo                    = "Info"
	TypeUpdate                  = "Update"
	TypeCommandResponse         = "CommandResponse"
	TypeCommandDownloadResponse = "CommandDownloadResponse"
	TypeCredentialResponse      = "CredentialResponse"
	TypeTokenResponse           = "TokenResponse"
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

type TokenResponse struct {
	Msg       string `json:"msg"`
	Token     string `json:"token"`
	PublicKey string `json:"publicKey"`
}

type CommandResponse struct {
	UpdateResponse
	SendCommandInfo
	Result string `json:"result"`
}

type CommandDownloadResponse struct {
	CommandResponse
	Filename string `json:"filename"`
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

func NewInfoResponse(message string) CustomHubResponse {
	return *NewCustomHubResponse(message, TypeInfo, StatusOk)
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

func NewCommandResponseInner(status int, result string, req SendCommandRequest) CommandResponse {
	updateResponse := NewUpdateResponseInner(req.ComponentId, status)
	inner := CommandResponse{
		UpdateResponse:  updateResponse,
		SendCommandInfo: req.SendCommandInfo,
		Result:          result,
	}
	return inner
}

func NewCommandResponse(status int, result string, req SendCommandRequest) CustomHubResponse {
	inner := NewCommandResponseInner(status, result, req)
	return *NewCustomHubResponse(inner, TypeCommandResponse, StatusOk)
}

func NewCommandDownloadResponse(status int, result string, filename string, req SendCommandRequest) CustomHubResponse {
	commandResponse := NewCommandResponseInner(status, result, req)
	inner := CommandDownloadResponse{
		CommandResponse: commandResponse,
		Filename:        filename,
	}
	return *NewCustomHubResponse(inner, TypeCommandDownloadResponse, StatusOk)
}

func NewCredentialResponse(message string) CustomHubResponse {
	return *NewCustomHubResponse(message, TypeCredentialResponse, StatusError)
}

func NewTokenResponse(message string, token string, publicKey string) CustomHubResponse {
	inner := TokenResponse{message, token, publicKey}
	return *NewCustomHubResponse(inner, TypeTokenResponse, StatusError)
}
