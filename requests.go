package main

import "encoding/json"

type SendCommandInfo struct {
	CommandId string `json:"commandId"`
	Command   string `json:"command"`
}

type SendCommandRequest struct {
	SendCommandInfo
	ComponentId string `json:"componentId"`
}

type SendVaultCredentialsRequest struct {
	Host  string `json:"host"`
	Token string `json:"token"`
	Path  string `json:"path"`
	Key   string `json:"key"`
}

func ReceiveRequest[R SendVaultCredentialsRequest | SendCommandRequest](message string) R {
	var req R
	_ = json.Unmarshal([]byte(message), &req)
	return req
}
