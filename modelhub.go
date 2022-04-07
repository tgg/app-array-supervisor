package main

import (
	"encoding/json"
	"fmt"
	"github.com/tgg/app-array-supervisor/model"
	"log"
)

type ModelHub struct {
	CustomHub
}

func NewModelHub() *ModelHub {
	hub := &ModelHub{}
	getAppArrayContext().SetModelHub(hub)
	return hub
}

func (h *ModelHub) SendModel(message string) {
	log.Printf("%s sent model: %s\n", h.ConnectionID(), message)
	var resp CustomHubResponse

	var a model.Application
	err := json.Unmarshal([]byte(message), &a)
	if err != nil {
		resp = NewErrorResponse(fmt.Sprintf("Incorrect model %s", err))
		h.SendResponseCaller(resp, "sendModelResponse")
	} else {
		ok, msg, ids := saveModel(a)
		if ok {
			var resp CustomHubResponse
			if msg != "" {
				resp = NewExistingModelResponse(a.Id, ids, msg)
			} else {
				resp = NewNewModelResponse(a.Id, ids)
				h.UpdateClients(resp, "newModelReceived")
			}
			h.SendResponseCaller(resp, "sendModelResponse")
		} else {
			h.SendResponseCaller(NewErrorResponse(msg), "sendModelResponse")
		}
	}
}

func saveModel(app model.Application) (bool, string, []string) {
	c := getAppArrayContext()
	if foundApp, ok := c.Models()[app.FormattedId()]; !ok {
		res := createHubs(app)
		if len(res) != 0 {
			c.Models()[app.FormattedId()] = app
			return true, "", res
		} else {
			return false, fmt.Sprintf("Model received but no signalr endpoint created : no environments provided for %s", app.Id), []string{}
		}
	} else {
		return true, fmt.Sprintf("Trying to register new model but model with id %s already exists", foundApp.Id), c.Paths()[foundApp.Id]
	}
}

func createHubs(app model.Application) []string {
	c := getAppArrayContext()
	var res []string
	for _, env := range app.Environments {
		hub := NewAppArrayHub(app, env)
		c.AppHubs()[env.FormattedId()] = hub
		c.Router().RegisterSignalRRoute(fmt.Sprintf("/%s/%s", app.FormattedId(), env.FormattedId()), hub)
		c.Paths()[app.FormattedId()] = append(c.Paths()[app.FormattedId()], hub.GetPath())
		res = append(res, hub.GetPath())
	}
	return res
}
