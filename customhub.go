package main

import (
	"encoding/json"
	"github.com/philippseith/signalr"
	"log"
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
	log.Printf("%s is connected with ID %s on : %s\n", h.RemoteAddr(), h.ConnectionID(), h.path)
}

func (h *CustomHub) SendResponseCaller(response CustomHubResponse, target string) {
	b, _ := json.Marshal(&response)
	h.Clients().Caller().Send(target, string(b))
}

func (h *CustomHub) UpdateClients(response CustomHubResponse, target string) {
	b, _ := json.Marshal(&response)
	h.Clients().All().Send(target, string(b))
}
