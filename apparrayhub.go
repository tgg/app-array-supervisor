package main

import (
	"fmt"
	"github.com/philippseith/signalr"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
)

type AppArrayHub struct {
	signalr.Hub
	sshClient *ssh.Client
}

func (h *AppArrayHub) OnConnected(string) {
	fmt.Printf("%s is connected\n", h.ConnectionID())
}

func (h *AppArrayHub) Coucou(message string) {

}

func (h *AppArrayHub) SendCommand(message string) {
	fmt.Printf("%s sent: %s\n", h.ConnectionID(), message)

	ctx := getAppArrayContext()
	fmt.Printf("%s\n", ctx.models[0])

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

	h.Clients().Client(h.ConnectionID()).Send("statusUpdated", buf.String())
	h.Clients().All().Send("statusUpdated", h.ConnectionID()+" sent "+message)
}
