package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sync"
)

type AuthenticationManagerInterface interface {
	CreateToken(connectionID string, path string) (string, error)
	ClientExists(token string)
	GetCredentials(login string) (string, bool)
	AddProvider(AuthenticationProviderInterface)
}

type AuthenticationManager struct {
	authManagers []AuthenticationProviderInterface
	clients      map[string]string
	cm           sync.Mutex
}

func NewAuthenticationManager() AuthenticationManagerInterface {
	return &AuthenticationManager{
		clients: map[string]string{},
	}
}

func (am *AuthenticationManager) CreateToken(connectionID string, path string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s/%s", connectionID, path)), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}
	hasher := md5.New()
	hasher.Write(hash)
	str := hex.EncodeToString(hasher.Sum(nil))
	return str, nil
}

func (am *AuthenticationManager) ClientExists(token string) {

}

func (am *AuthenticationManager) GetCredentials(login string) (string, bool) {
	pwd := ""
	ok := false
	for _, ap := range am.authManagers {
		pwd, ok = ap.GetCredentials(login)
		if ok {
			return pwd, ok
		}
	}
	return pwd, ok
}

func (am *AuthenticationManager) AddProvider(ap AuthenticationProviderInterface) {
	am.cm.Lock()
	defer am.cm.Unlock()
	am.authManagers = append(am.authManagers, ap)
}
