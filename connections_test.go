package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/tgg/app-array-supervisor/model"
	"testing"
)

func TestCreateSshClientsForApplication(t *testing.T) {
	c := getAppArrayContext()
	c.AddAuthManagers(&AuthenticationManager{map[string]string{"stephanes": "test"}})
	app := model.Application{
		Environments: []model.Environment{
			{
				model.Context{
					"Copier": map[string]string{"host": "localhost", "login": "stephanes", "output": "file:///tmp/my.tgz"},
					"Zipper": map[string]string{"host": "localhost", "login": "stephanes", "output": "file:///tmp/my.tgz"},
				},
				"This environment",
			},
		},
	}
	clients := CreateSshClientsForApplication(app.Environments[0])
	assert.EqualValues(t, 2, len(clients))
	assert.EqualValues(t, 1, len(getAppArrayContext().GetClientHosts()))
	for _, v := range getAppArrayContext().GetClientHosts() {
		v.Close()
	}
}

func TestCreateSshClientsForApplicationWithVault(t *testing.T) {
	c := getAppArrayContext()
	c.AddAuthManagers(NewVaultAuthenticationManager("http://localhost:8200", "dev-root-token", "kv", "password"))
	app := model.Application{
		Environments: []model.Environment{
			{
				model.Context{
					"Copier": map[string]string{"host": "localhost", "login": "stephanes", "output": "file:///tmp/my.tgz"},
					"Zipper": map[string]string{"host": "localhost", "login": "stephanes", "output": "file:///tmp/my.tgz"},
				},
				"This environment",
			},
		},
	}
	clients := CreateSshClientsForApplication(app.Environments[0])
	assert.EqualValues(t, 2, len(clients))
	assert.EqualValues(t, 1, len(getAppArrayContext().GetClientHosts()))
	for _, v := range getAppArrayContext().GetClientHosts() {
		v.Close()
	}
}
