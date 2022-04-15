package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tebeka/atexit"
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
	router.AddStaticFiles()
	router.RegisterSignalRRoute("/model", NewModelHub())
	log.Println("SignalR server created")

	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	log.Println("Cors allows all")
	handler := c.Handler(router)
	log.Println("Routes registered")

	getAppArrayContext().SetRouter(router)
	log.Println("Context configured")

	USE_SSL := os.Getenv("USE_SSL")
	listeningHost := "0.0.0.0:" + os.Getenv("LISTEN_PORT")
	if USE_SSL == "Y" {
		serverCRT := os.Getenv("SERVER_CRT")
		serverKEY := os.Getenv("SERVER_KEY")
		log.Printf("Listening for secured websocket connections on %s\n", listeningHost)
		if err := http.ListenAndServeTLS(listeningHost, serverCRT, serverKEY, handler); err != nil {
			log.Fatal("ListenAndServeTLS:", err)
		}
	} else {
		log.Printf("Listening for websocket connections on %s\n", listeningHost)
		if err := http.ListenAndServe(listeningHost, handler); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}

	atexit.Register(disconnectedClients)
	atexit.Exit(0)
}
