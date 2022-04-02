package main

import (
	"context"
	"github.com/tgg/app-array-supervisor/model"
	"sync"
)

const (
	AppArrayContextId = "AppArrayContextId"
)

type ConcurrentContext interface {
	Models() map[string]model.Application
	AppHubs() map[string]*AppArrayHub
	ModelHub() *ModelHub
	Router() *MuxRouterSignalR
	SetModelHub(*ModelHub)
	SetRouter(*MuxRouterSignalR)
}

type AppArrayContext struct {
	cm       sync.RWMutex
	models   map[string]model.Application
	appHub   map[string]*AppArrayHub
	modelHub *ModelHub
	router   *MuxRouterSignalR
}

var (
	globalCtx = context.WithValue(context.TODO(), AppArrayContextId,
		&AppArrayContext{
			models: map[string]model.Application{},
			appHub: map[string]*AppArrayHub{},
		})
)

func getAppArrayContext() ConcurrentContext {
	return globalCtx.Value(AppArrayContextId).(*AppArrayContext)
}

func (a *AppArrayContext) Models() map[string]model.Application {
	a.cm.RLock()
	defer a.cm.RUnlock()
	return a.models
}

func (a *AppArrayContext) AppHubs() map[string]*AppArrayHub {
	a.cm.RLock()
	defer a.cm.RUnlock()
	return a.appHub
}

func (a *AppArrayContext) ModelHub() *ModelHub {
	a.cm.RLock()
	defer a.cm.RUnlock()
	return a.modelHub
}

func (a *AppArrayContext) Router() *MuxRouterSignalR {
	a.cm.RLock()
	defer a.cm.RUnlock()
	return a.router
}

func (a *AppArrayContext) SetRouter(router *MuxRouterSignalR) {
	a.router = router
}

func (a *AppArrayContext) SetModelHub(hub *ModelHub) {
	a.modelHub = hub
}