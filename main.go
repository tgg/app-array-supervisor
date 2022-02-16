package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
	http.HandleFunc("/shell", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)

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
				log.Printf("Error running command: %v", err)
			}
		}
	})
	http.ListenAndServe(":"+os.Getenv("LISTEN_PORT"), nil)
}
