package model

import (
	"time"
)

type Status string

const (
	StatusUnknown  Status = "Unknown"
	StatusStarting        = "Starting"
	StatusStarted         = "Started"
	StatusRunning         = "Running"
	StatusStopping        = "Stopping"
	StatusStopped         = "Stopped"
)

type StatusTrail struct {
	Status Status    `json:"status"`
	At     time.Time `json:"at"`
	By     string    `json:"by"`
	On     string    `json:"on"`
}

type Command struct {
	Type  string   `json:"type"`
	Steps []string `json:"steps"`
}

type CommandMap map[string]Command

type TagMap map[string][]string

type Element struct {
	Name string  `json:"name"`
	Tags *TagMap `json:"tags,omitempty"`
}

type LiveElement struct {
	Element
	Commands CommandMap `json:"commands,omitempty"`
}

type BusinessObject string

type PortKind uint8

const (
	PortRead      PortKind = 1 << 1
	PortWrite              = 1 << 2
	PortReadWrite          = PortRead | PortWrite
)

type PortId string

type Port struct {
	Id       string         `json:"id"`
	Object   BusinessObject `json:"object,omitempty"`
	Kind     PortKind       `json:"kind"`
	Protocol string         `json:"protocol,omitempty"`
}

type Context map[string][]string

type Environment struct {
	Context
	Id string `json:"id"`
}

type Component struct {
	LiveElement
	Type     string   `json:"type"`
	Provides []Port   `json:"provides,omitempty"`
	Consumes []PortId `json:"consumes,omitempty"`
}

type Application struct {
	LiveElement
	Type         string        `json:"type"`
	Provides     []Port        `json:"provides,omitempty"`
	Consumes     []PortId      `json:"consumes,omitempty"`
	Components   []Component   `json:"components"`
	Environments []Environment `json:"environments,omitempty"`
}
