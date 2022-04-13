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
	paths []string
}

func (router *MuxRouterSignalR) HandleFunc(path string, f func(w http.ResponseWriter, r *http.Request)) {
	router.HandleFuncRoute(path, f)
}

func (router *MuxRouterSignalR) HandleFuncRoute(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	route := router.NewRoute().Path(path).HandlerFunc(f)
	router.paths = append(router.paths, path)
	log.Printf("Route %s registered\n", path)
	return route
}

func (router *MuxRouterSignalR) Handle(path string, handler http.Handler) {
	router.HandleRoute(path, handler)
}

func (router *MuxRouterSignalR) HandleRoute(path string, handler http.Handler) *mux.Route {
	route := router.NewRoute().Path(path).Handler(handler)
	router.paths = append(router.paths, path)
	log.Printf("Route %s registered\n", path)
	return route
}

func WithMuxRouter(router *MuxRouterSignalR) func() signalr.MappableRouter {
	return func() signalr.MappableRouter {
		return router
	}
}

func (router *MuxRouterSignalR) AddStaticFiles() {
	router.PathPrefix("/app-array/").Handler(http.StripPrefix("/app-array/", http.FileServer(http.Dir("./build/"))))
}

func (router *MuxRouterSignalR) AddHandledFunctions() {
	router.HandleFunc("/models/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		ctx := getAppArrayContext()
		id := vars["id"]
		if app, ok := ctx.Models()[id]; !ok {
			writer.WriteHeader(504)
		} else {
			if hub, ok := ctx.AppHubs()[id]; ok {
				hub.UpdateClientsWithID(NewMessageResponse(fmt.Sprintf("The address %s just requested the model %s", request.RemoteAddr, app.Id)), "statusUpdated", "router")
			}
			writer.WriteHeader(200)
			b, _ := json.Marshal(app)
			writer.Write(b)
		}
	})
	router.HandleFunc("/flush", func(w http.ResponseWriter, r *http.Request) {
		ctx := getAppArrayContext()
		i := 0
		for k := range ctx.Models() {
			delete(ctx.Models(), k)
			i++
		}
		l := 0
		for k := range ctx.AppHubs() {
			delete(ctx.AppHubs(), k)
			l++
		}
		w.Write([]byte(fmt.Sprintf("%d and %d deleted", i, l)))
	})
}

func (router *MuxRouterSignalR) RegisterSignalRRoute(path string, hub CustomHubInterface) {
	if !stringInSlice(path, router.paths) {
		hub.SetPath(path)
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
		server, _ := signalr.NewServer(context.TODO(),
			signalr.UseHub(hub),
			signalr.InsecureSkipVerify(true),
			signalr.KeepAliveInterval(2*time.Second),
			signalr.Logger(logger, false))
		server.MapHTTP(WithMuxRouter(router), path)
		log.Printf("SignalR route %s registered\n", path)
	} else {
		log.Printf("SignalR route %s not registered, it already exists\n", path)
	}
}
