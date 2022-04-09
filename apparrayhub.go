package main

import (
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
	sshClients  map[string]*ssh.Client
	credentials []Credential
	routines    []StatusRoutineInterface
	started     bool
	cm          sync.RWMutex
	app         model.Application
	env         model.Environment
}

func NewAppArrayHub(app model.Application, env model.Environment) *AppArrayHub {
	hub := &AppArrayHub{
		sshClients:  map[string]*ssh.Client{},
		credentials: GetHosts(env),
		app:         app,
		env:         env,
	}
	hub.routines = append(hub.routines, NewStatusRoutine(app, env, hub))
	return hub
}

func (h *AppArrayHub) InitializeConnections() {
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
	log.Printf("%s is connected on : %s\n", h.ConnectionID(), h.path)
	if len(h.credentials) > 0 && len(h.sshClients) != len(h.credentials) {
		h.SendResponseCaller(NewCredentialResponse("Credentials needed."), CredentialListener)
	}
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
	req := ReceiveRequest[SendCommandRequest](message)
	if req.Command != "" {
		if client, found := h.sshClients[req.ComponentId]; found {
			status, res := h.RunCommand(req.Command, client)
			h.SendResponseCaller(NewCommandResponse(status, res, req), CommandResultListener)
		} else {
			h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("No connection found for component %s", req.ComponentId)), StatusUpdateListener)
		}
	} else {
		h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("Bad request, cannot deserialize request : %s", message)), StatusUpdateListener)
	}
}

func (h *AppArrayHub) SendVaultCredentials(message string) {
	req := ReceiveRequest[SendVaultCredentialsRequest](message)
	c := getAppArrayContext()
	vam := NewVaultAuthenticationManager(req.Host, req.Token, req.Path, req.Key)
	if vam != nil {
		c.AddAuthManagers(vam)
		h.InitializeConnections()
		h.SendResponseCaller(NewInfoResponse("Credentials saved. Servers are accessible"), StatusUpdateListener)
	} else {
		h.SendResponseCaller(NewCredentialResponse("Vault credentials incorrect. Retry."), CredentialListener)
	}
}
