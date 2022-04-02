package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
)

type AppArrayHub struct {
	CustomHub
	sshClient *ssh.Client
}

func NewAppArrayHub(sshClient *ssh.Client) *AppArrayHub {
	hub := &AppArrayHub{
		sshClient: sshClient,
	}
	return hub
}

func (h *AppArrayHub) SendCommand(message string) {
	log.Printf("%s sent: %s\n", h.ConnectionID(), message)

	ctx := getAppArrayContext()
	log.Printf("%d\n", len(ctx.Models()))

	session, err := h.sshClient.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v\n", err)
	}

	if err := session.Run(string(message)); err != nil {
		log.Printf("Error running command: %v", err)
	}

	buf := new(strings.Builder)
	io.Copy(buf, stdout)

	h.SendResponseCaller(NewMessageResponse(buf.String()), "statusUpdated")
	h.UpdateClients(NewMessageResponse(fmt.Sprintf("%s sent %s", h.ConnectionID(), message)), "statusUpdated")
}
