package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func auth() (a ssh.AuthMethod, err error) {
	env := os.Getenv("REMOTE_SERVER_PK")
	if env != "" {
		key, err := ioutil.ReadFile(env)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
		return ssh.PublicKeys(signer), nil

	} else {
		return ssh.Password(os.Getenv("REMOTE_SERVER_PASSWORD")), nil
	}
}

// Simple http server exposing a websocket that will forward to ssh
func main() {
	auth, err := auth()
	if err != nil {
		log.Fatal("Failed to get authentication method: ", err)
	}
	config := &ssh.ClientConfig{
		User: os.Getenv("REMOTE_SERVER_USERNAME"),
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REMOTE_SERVER_HOST"), os.Getenv("REMOTE_SERVER_PORT")), config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	log.Printf("Connected to ssh server %v:%v\n", os.Getenv("REMOTE_SERVER_HOST"), os.Getenv("REMOTE_SERVER_PORT"))
	defer client.Close()

	router := &MuxRouterSignalR{
		Router: mux.NewRouter(),
	}
	router.AddHandledFunctions()
	router.RegisterSignalRRoute("/shell", NewAppArrayHub(client))
	router.RegisterSignalRRoute("/model", NewModelHub(client))
	log.Println("SignalR server created")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Requested-With", "X-Signalr-User-Agent"},
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
}
