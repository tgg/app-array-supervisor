package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type wsWriter struct {
	con     *websocket.Conn
	msgType int
}

func (w *wsWriter) Write(msg []byte) (n int, err error) {
	n = len(msg)
	err = w.con.WriteMessage(w.msgType, msg)
	return n, err
}

// Simple http server exposing a websocket that will forward to ssh
func main() {
	config := &ssh.ClientConfig{
		User: os.Args[1],
		Auth: []ssh.AuthMethod{
			ssh.Password(os.Args[2]),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", "localhost", "22"), config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	http.HandleFunc("/shell", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		if err != nil {
			log.Fatalln(err)
		}

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			session, err := client.NewSession()
			if err != nil {
				log.Fatal("Failed to create session: ", err)
			}
			defer session.Close()

			stdout, err := session.StdoutPipe()
			if err != nil {
				log.Fatalf("Unable to setup stdout for session: %v\n", err)
			}

			go func() {
				io.Copy(&wsWriter{conn, msgType}, stdout)
			}()

			if err := session.Run(string(msg)); err != nil {
				log.Fatal("Failed to run: " + err.Error())
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}
