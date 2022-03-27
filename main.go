package main

import (
	"context"
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/philippseith/signalr"
	"golang.org/x/crypto/ssh"
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
	fmt.Printf("Connected to ssh server %v:%v\n", os.Getenv("REMOTE_SERVER_HOST"), os.Getenv("REMOTE_SERVER_PORT"))
	defer client.Close()

	hub := AppArrayHub{
		sshClient: client,
	}

	server, _ := signalr.NewServer(context.TODO(),
		signalr.SimpleHubFactory(&hub),
		signalr.InsecureSkipVerify(true),
		//signalr.AllowOriginPatterns([]string{"localhost", "http://localhost:3000"}),
		signalr.KeepAliveInterval(2*time.Second),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stderr), true))

	router := http.NewServeMux()

	server.MapHTTP(signalr.WithHTTPServeMux(router), "/shell")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Requested-With", "X-Signalr-User-Agent"},
		Debug:            true,
	})
	handler := c.Handler(router)
	listeningHost := "localhost:" + os.Getenv("LISTEN_PORT")
	fmt.Printf("Listening for websocket connections on " + listeningHost)
	if err := http.ListenAndServe(listeningHost, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
