package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticationManager_CreateToken(t *testing.T) {
	am := NewAuthenticationManager()
	token, err := am.CreateToken("12q3sd1q32sdq=", "test")

	t.Log(token)

	assert.EqualValues(t, nil, err)
	assert.NotEqualValues(t, "", token)
}
