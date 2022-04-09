package main

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
	"net/url"
)

type AuthenticationManagerInterface interface {
	GetCredentials(login string) (string, bool)
}

const (
	VaultData = "data"
	VaultList = "metadata"
)

type AuthenticationManager struct {
	passwords map[string]string
}

func (am *AuthenticationManager) GetCredentials(login string) (string, bool) {
	if pwd, found := am.passwords[login]; found {
		return pwd, true
	} else {
		return "", false
	}
}

type VaultAuthenticationManager struct {
	client *vault.Client
	path   string
	key    string
}

func (am *VaultAuthenticationManager) GetCredentials(login string) (string, bool) {
	dataPath := fmt.Sprintf("%s/%s/%s", am.path, VaultData, login)
	secret, err := am.client.Logical().Read(dataPath)
	if err != nil || secret == nil {
		log.Printf("unable to read secret for %s: %v", login, err)
		return "", false
	}
	data, ok := secret.Data[VaultData].(map[string]interface{})
	if !ok {
		log.Printf("data type assertion failed: %T %#v", secret.Data[VaultData], secret.Data[VaultData])
		return "", false
	}

	value, ok := data[am.key].(string)
	if !ok {
		log.Printf("value type assertion failed: %T %#v", data[am.key], data[am.key])
		return "", false
	}

	return value, true
}

func NewVaultAuthenticationManager(host string, token string, path string, key string) *VaultAuthenticationManager {
	config := vault.DefaultConfig()
	config.Address = host
	client, err := vault.NewClient(config)
	if err != nil {
		log.Printf("unable to initialize Vault client: %v", err)
		return nil
	}
	client.SetToken(token)
	if foundData, err := client.Logical().List(fmt.Sprintf("%s/%s", path, VaultList)); err != nil {
		if apiErr, isResponseError := err.(*vault.ResponseError); isResponseError && apiErr.StatusCode == 403 {
			log.Printf("invalid token used for Vault %s: %v", host, err)
			return nil
		}
		if _, isUrlError := err.(*url.Error); isUrlError {
			log.Printf("incorrect url used for vault %s: %v", host, err)
			return nil
		}
	} else if foundData == nil {
		log.Printf("invalid path %s for Vault %s", path, host)
		return nil
	}
	return &VaultAuthenticationManager{client, path, key}
}
