package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/tgg/app-array-supervisor/model"
	"testing"
)

func TestCreateSshClientsForApplication(t *testing.T) {
	app := model.Application{
		Environments: []model.Environment{
			{
				model.Context{
					"Copier": map[string]string{"host": "localhost", "output": "file:///tmp/my.tgz"},
					"Zipper": map[string]string{"host": "localhost", "output": "file:///tmp/my.tgz"},
				},
				"This environment",
			},
		},
	}
	clients := CreateSshClientsForApplication(app)
	assert.EqualValues(t, 2, len(clients))
	assert.EqualValues(t, 1, len(getAppArrayContext().GetClientHosts()))
	for _, v := range getAppArrayContext().GetClientHosts() {
		v.Close()
	}
}
