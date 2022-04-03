package main

import (
	"fmt"
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

func Auth() (a ssh.AuthMethod, err error) {
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

func CreateSshClient(host string) *ssh.Client {
	auth, err := Auth()
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
	url := fmt.Sprintf("%s:%s", host, "22")
	client, err := ssh.Dial("tcp", url, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	log.Printf("Connected to ssh server %s\n", url)
	return client
}

//Returns map for component / ssh.client
func CreateSshClientsForApplication(app model.Application) map[string]*ssh.Client {
	res := map[string]*ssh.Client{}
	if len(app.Environments) != 0 {
		env := app.Environments[0]
		for componentId, v := range env.Context {
			if host, found := v["host"]; found {
				clientHosts := getAppArrayContext().GetClientHosts()
				var sshClient *ssh.Client
				if foundClient, found2 := clientHosts[host]; !found2 {
					sshClient = CreateSshClient(host)
					clientHosts[host] = sshClient
				} else {
					sshClient = foundClient
				}
				res[componentId] = sshClient
			}
		}
	}
	return res
}
