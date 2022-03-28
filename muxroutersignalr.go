package main

import (
	"github.com/gorilla/mux"
	"github.com/philippseith/signalr"
	"net/http"
)

type MuxRouterSignalR struct {
	*mux.Router
}

func (router *MuxRouterSignalR) HandleFunc(path string, f func(w http.ResponseWriter, r *http.Request)) {
	router.NewRoute().Path(path).HandlerFunc(f)
}

func (router *MuxRouterSignalR) Handle(path string, handler http.Handler) {
	router.NewRoute().Path(path).Handler(handler)
}

func WithMuxRouter(router *MuxRouterSignalR) func() signalr.MappableRouter {
	return func() signalr.MappableRouter {
		return router
	}
}
