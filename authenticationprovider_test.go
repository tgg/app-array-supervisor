package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewVaultAuthenticationManager(t *testing.T) {
	manager := NewVaultAuthenticationProvider("http://localhost:8200", "dev-root-token", "kv", "password")
	assert.True(t, manager != nil)
}

func TestNewVaultAuthenticationManagerWithIncorrectPath(t *testing.T) {
	manager := NewVaultAuthenticationProvider("http://localhost:8200", "dev-root-token", "kvv", "password")
	assert.True(t, manager == nil)
}

func TestNewVaultAuthenticationManagerWithIncorrectToken(t *testing.T) {
	manager := NewVaultAuthenticationProvider("http://localhost:8200", "false-token", "kv", "password")
	assert.True(t, manager == nil)
}

func TestNewVaultAuthenticationManagerWithIncorrectHost(t *testing.T) {
	manager := NewVaultAuthenticationProvider("http://localhost:8201", "dev-root-token", "kv", "password")
	assert.True(t, manager == nil)
}

func TestVaultAuthenticationManager_GetCredentials(t *testing.T) {
	manager := NewVaultAuthenticationProvider("http://localhost:8200", "dev-root-token", "kv", "password")
	pwd, ok := manager.GetCredentials("test")
	assert.True(t, ok)
	assert.EqualValues(t, "pouet", pwd)
}

func TestVaultAuthenticationManager_GetCredentialsUnknownLogin(t *testing.T) {
	manager := NewVaultAuthenticationProvider("http://localhost:8200", "dev-root-token", "kv", "secret")
	pwd, ok := manager.GetCredentials("testnotexist")
	assert.False(t, ok)
	assert.EqualValues(t, "", pwd)
}

func TestVaultAuthenticationManager_GetCredentialsIncorrectKey(t *testing.T) {
	manager := NewVaultAuthenticationProvider("http://localhost:8200", "dev-root-token", "kv", "secret")
	pwd, ok := manager.GetCredentials("test")
	assert.False(t, ok)
	assert.EqualValues(t, "", pwd)
}
