package main

import (
	"context"
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"sync"
)

const (
	AppArrayContextId = "AppArrayContextId"
)

type ConcurrentContext interface {
	Models() map[string]model.Application
	AppHubs() map[string]*AppArrayHub
	Paths() map[string][]string
	ModelHub() *ModelHub
	Router() *MuxRouterSignalR
	SetModelHub(*ModelHub)
	SetRouter(*MuxRouterSignalR)
	GetClientHosts() map[string]*ssh.Client
	GetAuthManagers() []AuthenticationManagerInterface
	AddAuthManagers(am AuthenticationManagerInterface)
}

type AppArrayContext struct {
	cm       sync.RWMutex
	models   map[string]model.Application
	paths    map[string][]string
	appHub   map[string]*AppArrayHub
	modelHub *ModelHub
	router   *MuxRouterSignalR

	//TODO: should be per client
	clientHost   map[string]*ssh.Client
	authManagers []AuthenticationManagerInterface
}

var (
	globalCtx = context.WithValue(context.TODO(), AppArrayContextId,
		&AppArrayContext{
			models:     map[string]model.Application{},
			paths:      map[string][]string{},
			appHub:     map[string]*AppArrayHub{},
			clientHost: map[string]*ssh.Client{},
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

func (a *AppArrayContext) Paths() map[string][]string {
	a.cm.RLock()
	defer a.cm.RUnlock()
	return a.paths
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

func (a *AppArrayContext) GetClientHosts() map[string]*ssh.Client {
	a.cm.RLock()
	defer a.cm.RUnlock()
	return a.clientHost
}

func (a *AppArrayContext) SetRouter(router *MuxRouterSignalR) {
	a.router = router
}

func (a *AppArrayContext) SetModelHub(hub *ModelHub) {
	a.modelHub = hub
}

func (a *AppArrayContext) AddAuthManagers(am AuthenticationManagerInterface) {
	a.authManagers = append(a.authManagers, am)
}

func (a *AppArrayContext) GetAuthManagers() []AuthenticationManagerInterface {
	return a.authManagers
}
