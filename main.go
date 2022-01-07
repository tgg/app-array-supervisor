package main

import (
	"os"

	"github.com/philippseith/signalr"
	"github.com/tgg/app-array-supervisor/remote"
)

type Supervisor struct {
	signalr.Hub
}

func (s *Supervisor) RunCommand(command string, args []string) {

}

func main() {
	remote.NewClient("localhost", "22", os.Args[1], os.Args[2])
}
