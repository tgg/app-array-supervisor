package main

import (
	"errors"
	"fmt"
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

func GetPassword(login string) (string, bool) {
	c := getAppArrayContext()
	return c.GetAuthManager().GetCredentials(login)
}

func Auth(login string) (a ssh.AuthMethod, err error) {
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
		pwd, ok := GetPassword(login)
		if ok {
			return ssh.Password(pwd), nil
		} else {
			return nil, errors.New(fmt.Sprintf("Password not found for %s", login))
		}
	}
}

func CreateSshClient(host string, login string) *ssh.Client {
	auth, err := Auth(login)
	if err != nil {
		log.Printf("Failed to get authentication method: %v", err)
		return nil
	}
	config := &ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	url := fmt.Sprintf("%s:%s", host, "22")
	client, err := ssh.Dial("tcp", url, config)
	if err != nil {
		log.Printf("Failed to dial: %v", err)
		return nil
	}
	log.Printf("Connected to ssh server %s\n", url)
	return client
}

type Credential struct {
	Component string
	Host      string
	Login     string
}

func GetHosts(env model.Environment) []Credential {
	var credentials []Credential
	for componentId, v := range env.Context {
		if host, foundHost := v["host"]; foundHost {
			if login, foundLogin := v["login"]; foundLogin {
				credentials = append(credentials, Credential{componentId, host, login})
			}
		}
	}
	return credentials
}

//Returns map for component / ssh.client
func CreateSshClientsForApplication(env model.Environment) map[string]*ssh.Client {
	res := map[string]*ssh.Client{}
	for _, cred := range GetHosts(env) {
		clientHosts := getAppArrayContext().GetClientHosts()
		var sshClient *ssh.Client
		if foundClient, found2 := clientHosts[cred.Host]; !found2 {
			sshClient = CreateSshClient(cred.Host, cred.Login)
			if sshClient != nil {
				clientHosts[cred.Host] = sshClient
				res[cred.Component] = sshClient
			}
		} else {
			sshClient = foundClient
			res[cred.Component] = sshClient
		}
	}
	return res
}
