package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tebeka/atexit"
	"log"
	"net/http"
	"os"
)

func disconnectedClients() {
	ctx := getAppArrayContext()
	for _, v := range ctx.GetClientHosts() {
		v.Close()
	}
}

// Simple http server exposing a websocket that will forward to ssh
func main() {
	router := &MuxRouterSignalR{
		Router: mux.NewRouter(),
	}
	router.AddHandledFunctions()
	router.RegisterSignalRRoute("/model", NewModelHub())
	log.Println("SignalR server created")

	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	handler := c.Handler(router)
	log.Println("Routes registered")

	getAppArrayContext().SetRouter(router)
	log.Println("Context configured")

	listeningHost := "localhost:" + os.Getenv("LISTEN_PORT")
	log.Printf("Listening for websocket connections on %s\n", listeningHost)
	if err := http.ListenAndServe(listeningHost, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	atexit.Register(disconnectedClients)
	atexit.Exit(0)
}
