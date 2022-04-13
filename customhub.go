package main

import (
	"encoding/json"
	"github.com/philippseith/signalr"
	"log"
)

const (
	StatusUpdateListener  = "statusUpdated"
	CommandResultListener = "getCommandResult"
	CredentialListener    = "onCredentialResponse"
	SendModelListener     = "sendModelResponse"
	ModelReceivedListener = "newModelReceived"
	TokenReceivedListener = "tokenReceived"
)

type CustomHubInterface interface {
	Initialize(hubContext signalr.HubContext)
	OnConnected(connectionID string)
	OnDisconnected(connectionID string)
	SetPath(path string)
	GetPath() string
}

type CustomHub struct {
	signalr.Hub
	path string
}

func (h *CustomHub) SetPath(path string) {
	h.path = path
}

func (h *CustomHub) GetPath() string {
	return h.path
}

func (h *CustomHub) OnConnected(string) {
	h.Groups().AddToGroup(h.path, h.ConnectionID())
	log.Printf("%s is connected on : %s\n", h.ConnectionID(), h.path)
}

func (h *CustomHub) OnDisconnected(string) {
	h.Groups().RemoveFromGroup(h.path, h.ConnectionID())
	log.Printf("%s is disconnected from : %s\n", h.ConnectionID(), h.path)
}

func (h *CustomHub) SendResponseCaller(response CustomHubResponse, target string) {
	log.Printf("Route %s : Sending response %d with message \"%s\" to %s\n", h.path, response.StatusCode, response.Msg, h.ConnectionID())
	b, _ := json.Marshal(&response)
	h.Clients().Caller().Send(target, string(b))
}

func (h *CustomHub) UpdateClientsWithID(response CustomHubResponse, target string, connectionID string) {
	log.Printf("Route %s : Update all clients with message \"%s\" from %s\n", h.path, response.Msg, connectionID)
	b, _ := json.Marshal(&response)
	h.Clients().Group(h.path).Send(target, string(b))
}

func (h *CustomHub) UpdateClients(response CustomHubResponse, target string) {
	h.UpdateClientsWithID(response, target, h.ConnectionID())
}
