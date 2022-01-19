package model

import (
	"encoding/json"
	"errors"
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

// Because we want this model to follow TypeScript we need to
// refine marshaling / unmarshalling
func (e *Environment) MarshalJSON() ([]byte, error) {
	// We need to encode the Context, then add Id
	copy := make(Context)
	for key, value := range e.Context {
		copy[key] = value
	}
	copy["id"] = []string{e.Id}
	return json.Marshal(copy)
}

func (e *Environment) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &e.Context)

	if err != nil {
		return err
	}

	if ids, ok := e.Context["id"]; !ok || len(ids) != 1 {
		return errors.New(`Missing or incorrect "id" field`)

	} else {
		e.Id = ids[0]
		delete(e.Context, "id")
		return nil
	}
}
