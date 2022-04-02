package main

import (
	"context"
	"encoding/json"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/philippseith/signalr"
	"log"
	"net/http"
	"os"
	"time"
)

type MuxRouterSignalR struct {
	*mux.Router
}

func (router *MuxRouterSignalR) HandleFunc(path string, f func(w http.ResponseWriter, r *http.Request)) {
	router.NewRoute().Path(path).HandlerFunc(f)
	log.Printf("Route %s registered\n", path)
}

func (router *MuxRouterSignalR) Handle(path string, handler http.Handler) {
	router.NewRoute().Path(path).Handler(handler)
	log.Printf("Route %s registered\n", path)
}

func WithMuxRouter(router *MuxRouterSignalR) func() signalr.MappableRouter {
	return func() signalr.MappableRouter {
		return router
	}
}

func (router *MuxRouterSignalR) AddHandledFunctions() {
	router.HandleFunc("/models/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		ctx := getAppArrayContext()
		if app, ok := ctx.Models()[vars["id"]]; !ok {
			writer.WriteHeader(504)
		} else {
			ctx.ModelHub().Clients().All().Send("statusUpdated", fmt.Sprintf("The adress %s just requested the model %s", request.RemoteAddr, app.Id))
			writer.WriteHeader(200)
			b, _ := json.Marshal(app)
			writer.Write(b)
		}
	})
}

func (router *MuxRouterSignalR) RegisterSignalRRoute(path string, hub CustomHubInterface) {
	hub.SetPath(path)
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	server, _ := signalr.NewServer(context.TODO(),
		signalr.UseHub(hub),
		signalr.InsecureSkipVerify(true),
		signalr.KeepAliveInterval(2*time.Second),
		signalr.Logger(logger, false))
	server.MapHTTP(WithMuxRouter(router), path)
	log.Printf("SignalR route %s registered\n", path)
}
