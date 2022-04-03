package main

import (
	"encoding/json"
	"fmt"
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
)

type AppArrayHub struct {
	CustomHub
	sshClients map[string]*ssh.Client
}

func NewAppArrayHub(app model.Application) *AppArrayHub {
	hub := &AppArrayHub{
		sshClients: CreateSshClientsForApplication(app),
	}
	if len(app.Environments) != 0 {
		hub.routines = append(hub.routines, NewStatusRoutine(app, app.Environments[0], hub))
	}
	return hub
}

type SendCommandRequest struct {
	Command   string `json:"command"`
	Component string `json:"id"`
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
		if client, found := h.sshClients[req.Component]; found {
			_, res := h.RunCommand(req.Command, client)
			h.SendResponseCaller(NewMessageResponse(res), "getCommandResult")
		} else {
			h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("No connection found for component %s", req.Component)), "statusUpdated")
		}
	} else {
		h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("Bad request, cannot deserialize request : %s", message)), "statusUpdated")
	}
}
