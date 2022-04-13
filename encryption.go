package main

import (
	"crypto/rand"
	"crypto/rsa"
	x509 "crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
)

type EncryptionInterface interface {
	Encrypt(string, string) ([]byte, error)
	EncryptToBase64(string, string) (string, error)
	Decrypt([]byte) (string, error)
	DecryptFromBase64(string) (string, error)
	GetPublicKey() string
}

type RSAEncryption struct {
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	publicKeyPem string
}

func (e *RSAEncryption) Encrypt(msg string, key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return []byte{}, errors.New("cannot parse pem")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Printf("Cannot parse publicKey : %v", err)
		return []byte{}, err
	}
	encryptedMsg, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), []byte(msg))
	if err != nil {
		log.Printf("Cannot encrypt message : %v", err)
		return []byte{}, err
	}
	return encryptedMsg, err
}

func (e *RSAEncryption) EncryptToBase64(msg string, key string) (string, error) {
	encrypted, err := e.Encrypt(msg, key)
	if err != nil {
		log.Printf("Cannot encrypt message : %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), err
}

func (e *RSAEncryption) Decrypt(msg []byte) (string, error) {
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, e.privateKey, msg)
	if err != nil {
		log.Printf("Cannot decrypt message : %v", err)
		return "", err
	}
	return string(decryptedBytes), nil
}

func (e *RSAEncryption) DecryptFromBase64(msg string) (string, error) {
	str, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		log.Printf("Cannot decrypt base64 message : %v", err)
		return "", err
	}
	return e.Decrypt(str)
}

func (e *RSAEncryption) GetPublicKey() string {
	return e.publicKeyPem
}

func NewRSAEncryption() EncryptionInterface {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// The public key is a part of the *rsa.PrivateKey struct
	publicKey := privateKey.PublicKey
	pubkey_bytes, _ := x509.MarshalPKIXPublicKey(&publicKey)
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)
	log.Println("RSA Encryption created")
	return &RSAEncryption{privateKey: privateKey, publicKey: &publicKey, publicKeyPem: string(pubkey_pem)}
}
