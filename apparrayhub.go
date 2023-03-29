package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/tgg/app-array-supervisor/model"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
	"sync"
	"os/exec"
	"runtime"
)

type AppArrayHub struct {
	CustomHub
	sshClients  map[string]*ssh.Client
	credentials []Credential
	routines    []StatusRoutineInterface
	started     bool
	cm          sync.RWMutex
	app         model.Application
	env         model.Environment
}

func NewAppArrayHub(app model.Application, env model.Environment) *AppArrayHub {
	hub := &AppArrayHub{
		sshClients:  map[string]*ssh.Client{},
		credentials: GetHosts(env),
		app:         app,
		env:         env,
	}
	hub.routines = append(hub.routines, NewStatusRoutine(app, env, hub))
	return hub
}

func (h *AppArrayHub) InitializeConnections() {
	h.cm.Lock()
	defer h.cm.Unlock()
	if !h.started {
		h.sshClients = CreateSshClientsForApplication(h.env)
		for _, routine := range h.routines {
			routine.Run()
		}
		h.started = true
	}
}

func (h *AppArrayHub) OnConnected(string) {
	h.Groups().AddToGroup(h.path, h.ConnectionID())
	log.Printf("%s is connected on : %s\n", h.ConnectionID(), h.path)
	if len(h.credentials) > 0 && len(h.sshClients) != len(h.credentials) {
		sshClients := CreateSshClientsForApplication(h.env)
		if len(sshClients) != len(h.credentials) {
			h.SendResponseCaller(NewCredentialResponse("Credentials needed."), CredentialListener)
		} else {
			h.InitializeConnections()
		}
	}

}

func (h *AppArrayHub) RunCommand(cmd string, client *ssh.Client) (int, string) {
	session, err := client.NewSession()
	res := 0
	if err != nil {
		log.Printf("Failed to create session: %v\n", err)
		return 1, ""
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("Unable to setup stdout for session: %v\n", err)
		return 1, ""
	}

	if errCmd := session.Run(cmd); errCmd != nil {
		log.Printf("Error running command: %v", errCmd)
		if sshErr, isExitError := errCmd.(*ssh.ExitError); isExitError {
			res = sshErr.ExitStatus()
		}
	}

	buf := new(strings.Builder)
	io.Copy(buf, stdout)
	return res, buf.String()
}
// open browser with url
func (h *AppArrayHub) OpenUrl(url string, client *ssh.Client) (int, string) {
	session, err := client.NewSession()
	//var err error
	res := 0
	if err != nil {
		log.Printf("Failed to create session: %v\n", err)
		return 1,""
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("Unable to setup stdout for session: %v\n", err)
		return 1,""
	}

	switch runtime.GOOS {
	case "linux":
		if err := session.Run("xdg-open " + url); err != nil {
			log.Printf("Error running command: %v", err)
			if sshErr, isExitError := err.(*ssh.ExitError); isExitError {
				res = sshErr.ExitStatus()
			}
		}
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		if err := session.Run("rundll32, url.dll,FileProtocolHandler " + url); err != nil {
			log.Printf("Error running command: %v", err)
			if sshErr, isExitError := err.(*ssh.ExitError); isExitError {
				res = sshErr.ExitStatus()
			}
		}
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		if err := session.Run("open " + url); err != nil {
			log.Printf("Error running command: %v", err)
			if sshErr, isExitError := err.(*ssh.ExitError); isExitError {
				res = sshErr.ExitStatus()
			}
		}
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		res = 1
	}
	
	buf := new(strings.Builder)
	io.Copy(buf, stdout)
	return res, buf.String()
}

// Open a terminal
func (h *AppArrayHub) OpenTerminal(url string, client *ssh.Client) (int,string) {
	session, err := client.NewSession()
	res := 0

	if err != nil {
		log.Printf("Failed to create session: %v\n", err)
		return 1,""
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("Unable to setup stdout for session: %v\n", err)
		return 1,""
	}

	switch runtime.GOOS {
	case "linux":
		if err := session.Run("cd " + url); err != nil {
			log.Printf("Error running command: %v", err)
			if sshErr, isExitError := err.(*ssh.ExitError); isExitError {
				res = sshErr.ExitStatus()
			}
		}
		err = exec.Command("cd",url).Start()
	case "windows":
		if err := session.Run("cmd.exe /c start cd " + url); err != nil {
			log.Printf("Error running command: %v", err)
			if sshErr, isExitError := err.(*ssh.ExitError); isExitError {
				res = sshErr.ExitStatus()
			}
		}
		err = exec.Command("cmd.exe","/c","start","cd",url).Start()
	default:
		err = fmt.Errorf("unsupported command-line")
	}

	if err != nil {
		log.Fatal(err)
		res = 1
	}
	for i := 0; i < 100; i++{
		log.Println(i)
	}

	buf := new(strings.Builder)
	io.Copy(buf, stdout)
	return res,buf.String()
}

func (h *AppArrayHub) DownloadFile(path string, client *ssh.Client) (int, string, string) {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Printf("Failed to open sftp session : %v", err)
		return 1, "", ""
	}

	file, err := sftpClient.Open(path)
	if err != nil {
		log.Printf("Could not open file %s : %v", path, err)
		return 1, "", ""
	}

	buf := new(strings.Builder)
	file.WriteTo(buf)
	return 0, buf.String(), file.Name()
}

func (h *AppArrayHub) SendCommand(message string) {
	log.Printf("Route %s : %s sent: %s\n", h.path, h.ConnectionID(), message)
	req := ReceiveRequest[SendCommandRequest](message)
	if req.Command != "" {
		if client, found := h.sshClients[req.ComponentId]; found {
			var res, filename string
			var status int
			if req.CommandId == model.CommandDownload {
				status, res, filename = h.DownloadFile(req.Command, client)
				h.SendResponseCaller(NewCommandDownloadResponse(status, res, filename, req), CommandResultListener)
			} else if req.CommandId == model.CommandWebsite {
				status, res = h.OpenUrl(req.Command,client)
				h.SendResponseCaller(NewCommandResponse(status, res, req), CommandResultListener)
			} else if req.CommandId == model.CommandTerminal {
				status, res = h.OpenTerminal(req.Command,client)
				h.SendResponseCaller(NewCommandResponse(status, res, req), CommandResultListener)
			} else {
				status, res = h.RunCommand(req.Command, client)
				h.SendResponseCaller(NewCommandResponse(status, res, req), CommandResultListener)
			}
		} else {
			h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("No connection found for component %s", req.ComponentId)), StatusUpdateListener)
		}
	} else {
		h.SendResponseCaller(NewErrorResponse(fmt.Sprintf("Bad request, cannot deserialize request : %s", message)), StatusUpdateListener)
	}
}

func (h *AppArrayHub) RequestToken() {
	c := getAppArrayContext()
	pem := c.GetEncryption().GetPublicKey()
	token, _ := c.GetAuthManager().CreateToken(h.ConnectionID(), h.GetPath())
	h.SendResponseCaller(NewTokenResponse("", token, pem), TokenReceivedListener)
}

func (h *AppArrayHub) SendTextCredentials(message string) {
	c := getAppArrayContext()
	decrypted, err := c.GetEncryption().DecryptFromBase64(message)
	if err != nil {
		h.SendResponseCaller(NewErrorResponse("Cannot decrypt sent message"), StatusUpdateListener)
		return
	} else {
		req := ReceiveRequest[SendTextCredentialsRequest](decrypted)
		vap := NewAuthenticationProvider(req.Credentials)
		c.GetAuthManager().AddProvider(vap)
		h.InitializeConnections()
		h.SendResponseCaller(NewInfoResponse("Credentials saved. Servers are accessible"), StatusUpdateListener)
	}
}

func (h *AppArrayHub) SendVaultCredentials(message string) {
	c := getAppArrayContext()
	decrypted, err := c.GetEncryption().DecryptFromBase64(message)
	if err != nil {
		h.SendResponseCaller(NewErrorResponse("Cannot decrypt sent message"), StatusUpdateListener)
		return
	} else {
		req := ReceiveRequest[SendVaultCredentialsRequest](decrypted)
		vap := NewVaultAuthenticationProvider(req.Host, req.Token, req.Path, req.Key)
		if vap != nil {
			c.GetAuthManager().AddProvider(vap)
			h.InitializeConnections()
			h.SendResponseCaller(NewInfoResponse("Credentials saved. Servers are accessible"), StatusUpdateListener)
		} else {
			h.SendResponseCaller(NewCredentialResponse("Vault credentials incorrect. Retry."), CredentialListener)
		}
	}
}
