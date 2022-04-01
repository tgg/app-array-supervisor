package main

import (
	"github.com/philippseith/signalr"
	"log"
)

type ModelHub struct {
	signalr.Hub
}

func (h *ModelHub) OnConnected(string) {
	log.Printf("%s is connected\n", h.ConnectionID())
}
