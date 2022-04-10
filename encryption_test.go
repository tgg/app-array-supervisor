package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecryptRSA(t *testing.T) {
	encryption := NewRSAEncryption()
	msg := "Coucou"
	key := encryption.GetPublicKey()
	t.Log(key)

	encodedMsg, _ := encryption.Encrypt(msg, key)
	t.Log(encodedMsg)

	decodedMsg, _ := encryption.Decrypt(encodedMsg)

	assert.EqualValues(t, msg, decodedMsg)
}

func TestEncryptDecryptRSABase64(t *testing.T) {
	encryption := NewRSAEncryption()
	msg := "Coucou"
	key := encryption.GetPublicKey()
	t.Log(key)

	encodedMsg, _ := encryption.EncryptToBase64(msg, key)
	t.Log(encodedMsg)

	decodedMsg, _ := encryption.DecryptFromBase64(encodedMsg)

	assert.EqualValues(t, msg, decodedMsg)
}
