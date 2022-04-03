package main

import (
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

type StatusRoutineInterface interface {
	Run()
}

type StatusComponentRoutine struct {
	c         model.Component
	cmd       model.Command
	hub       *AppArrayHub
	sshClient *ssh.Client
}

type StatusRoutine struct {
	app model.Application
	env model.Environment
	hub *AppArrayHub
}

func NewStatusRoutine(app model.Application, env model.Environment, hub *AppArrayHub) StatusRoutineInterface {
	return &StatusRoutine{app, env, hub}
}

func (s *StatusComponentRoutine) Run() {
	log.Printf("Status routine for route %s, component %s is running\n", s.hub.path, s.c.Id)
	for {
		time.Sleep(5 * time.Second)
		res := s.hub.RunCommand(s.cmd.Steps[0], s.sshClient)
		s.hub.UpdateClientsWithID(NewMessageResponse(res), "statusUpdated", "routine")
	}
	return
}

func (s *StatusRoutine) Run() {
	for id, env := range s.env.Context {
		if _, ok := env["host"]; ok {
			if sshClient, found := s.hub.sshClients[id]; found {
				for _, comp := range s.app.Components {
					if cmd, foundCmd := comp.Commands["status"]; comp.Id == id && foundCmd {
						routine := StatusComponentRoutine{comp, cmd, s.hub, sshClient}
						go routine.Run()
					}
				}
			}
		}
	}
}
