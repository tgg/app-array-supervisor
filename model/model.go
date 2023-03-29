package model

import (
	"encoding/json"
	"errors"
	"strings"
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

const (
	CommandStart    = "start"
	CommandStop     = "stop"
	CommandStatus   = "status"
	CommandDownload = "download"
	CommandWebsite  = "website"
	CommandTerminal = "terminal"
)

type Command struct {
	Type  string   `json:"type"`
	Steps []string `json:"steps"`
}

type CommandMap map[string]Command

type TagMap map[string]string

type Element struct {
	Id   string  `json:"id"`
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

type Context map[string]map[string]string
type ContextExport map[string]any

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
	Environments []Environment `json:"environments"`
}

func formattedId(id string) string {
	return strings.Replace(id, " ", "_", -1)
}

func (l *LiveElement) FormattedId() string {
	return formattedId(l.Id)
}

func (e *Environment) FormattedId() string {
	return formattedId(e.Id)
}

// Because we want this model to follow TypeScript we need to
// refine marshaling / unmarshalling
func (e *Environment) MarshalJSON() ([]byte, error) {
	// We need to encode the Context, then add Id
	copy := make(ContextExport)
	// Warning: we are not doing a deep copy here.
	for key, value := range e.Context {
		copy[key] = value
	}
	copy["id"] = e.Id
	return json.Marshal(copy)
}

func (e *Environment) UnmarshalJSON(data []byte) error {
	contextExport := ContextExport{}
	e.Context = Context{}
	err := json.Unmarshal(data, &contextExport)

	if err != nil {
		return err
	}

	if id, ok := contextExport["id"]; !ok {
		return errors.New(`Missing or incorrect "id" field`)
	} else {
		e.Id = id.(string)
		delete(contextExport, "id")
		for key, value := range contextExport {
			valueMap := value.(map[string]any)
			e.Context[key] = map[string]string{}
			for k, v := range valueMap {
				e.Context[key][k] = v.(string)
			}
		}
		return nil
	}
	return nil
}
