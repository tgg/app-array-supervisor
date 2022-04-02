package main

import (
	"encoding/json"
	"fmt"
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"log"
)

type ModelHub struct {
	CustomHub
	sshClient *ssh.Client
}

func NewModelHub(sshClient *ssh.Client) *ModelHub {
	hub := &ModelHub{
		sshClient: sshClient,
	}
	getAppArrayContext().SetModelHub(hub)
	return hub
}

func (h *ModelHub) SendModel(message string) {
	log.Printf("%s sent model: %s\n", h.ConnectionID(), message)
	var resp CustomHubResponse

	var a model.Application
	err := json.Unmarshal([]byte(message), &a)
	if err != nil {
		resp = NewErrorResponse(fmt.Sprintf("Incorrect model %s", err))
		h.SendResponseCaller(resp, "sendModelResponse")
	} else {
		ok, msg := SaveModel(a, h.sshClient)
		if ok {
			resp = NewNewModelResponse(a.Id, getAppArrayContext().AppHubs()[a.Id].GetPath())
			h.SendResponseCaller(resp, "sendModelResponse")
			h.UpdateClients(resp, "newModelReceived")
		} else {
			h.SendResponseCaller(NewErrorResponse(msg), "sendModelResponse")
		}
	}
}

func SaveModel(app model.Application, client *ssh.Client) (bool, string) {
	c := getAppArrayContext()
	if foundApp, ok := c.Models()[app.Id]; !ok {
		c.Models()[app.Id] = app
		hub := NewAppArrayHub(client)
		if foundHub, ok2 := c.AppHubs()[app.Id]; !ok2 {
			c.AppHubs()[app.Id] = hub
			c.Router().RegisterSignalRRoute(fmt.Sprintf("/%s", app.Id), hub)
			return true, ""
		} else {
			return false, fmt.Sprintf("Trying to register new hub for id %s but it already exists", foundHub.GetPath())
		}
	} else {
		return false, fmt.Sprintf("Trying to register new model but model with id %s already exists", foundApp.Id)
	}
}
