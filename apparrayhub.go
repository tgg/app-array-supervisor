package main

import (
	"encoding/json"
	"fmt"
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
	"sync"
)

type AppArrayHub struct {
	CustomHub
	sshClients map[string]*ssh.Client
	routines   []StatusRoutineInterface
	started    bool
	cm         sync.RWMutex
	app        model.Application
	env        model.Environment
}

func NewAppArrayHub(app model.Application, env model.Environment) *AppArrayHub {
	hub := &AppArrayHub{
		sshClients: map[string]*ssh.Client{},
		app:        app,
		env:        env,
	}
	if len(app.Environments) != 0 {
		hub.routines = append(hub.routines, NewStatusRoutine(app, env, hub))
	}
	return hub
}

func (h *AppArrayHub) OnFirstConnected() {
	h.cm.Lock()
	defer h.cm.Unlock()
	if !h.started {
		h.sshClients = CreateSshClientsForApplication(h.env)
		for _, routine := range h.routines {
			routine.Run()
		}
		h.started = true
	}
}

func (h *AppArrayHub) OnConnected(string) {
	h.Groups().AddToGroup(h.path, h.ConnectionID())
	h.OnFirstConnected()
	log.Printf("%s is connected on : %s\n", h.ConnectionID(), h.path)
}

type SendCommandInfo struct {
	CommandId string `json:"commandId"`
	Command   string `json:"command"`
}

type SendCommandRequest struct {
	SendCommandInfo
	ComponentId string `json:"componentId"`
}

func ReceiveSendCommandRequest(message string) SendCommandRequest {
	var req SendCommandRequest
	err := json.Unmarshal([]byte(message), &req)
	if err != nil {
		req = SendCommandRequest{}
	}
	return req
}

func (h *AppArrayHub) RunCommand(cmd string, client *ssh.Client) (int, string) {
	session, err := client.NewSession()
	res := 0
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v\n", err)
	}

	if errCmd := session.Run(cmd); errCmd != nil {
		log.Printf("Error running command: %v", errCmd)
		if sshErr, isExitError := errCmd.(*ssh.ExitError); isExitError {
			res = sshErr.ExitStatus()
		}
	}

	buf := new(strings.Builder)
	io.Copy(buf, stdout)
	return res, buf.String()
}

func (h *AppArrayHub) SendCommand(message string) {
	log.Printf("Route %s : %s sent: %s\n", h.path, h.ConnectionID(), message)
	req := ReceiveSendCommandRequest(message)
	if req.Command != "" {
		if client, found := h.sshClients[req.ComponentId]; found {
			status, res := h.RunCommand(req.Command, client)
			h.SendResponseCaller(NewCommandResponse(status, res, req), "getCommandResult")
		} else {
			h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("No connection found for component %s", req.ComponentId)), "statusUpdated")
		}
	} else {
		h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("Bad request, cannot deserialize request : %s", message)), "statusUpdated")
	}
}
