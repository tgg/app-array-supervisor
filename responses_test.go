package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNewModelResponseMsg(t *testing.T) {
	resp := NewNewModelResponse("Demo", "/Demo")
	b, _ := json.Marshal(&resp.Msg)
	msg := `{"id":"Demo","path":"/Demo"}`
	if !assert.EqualValues(t, msg, string(b)) {
		t.Errorf(`Deserialisation of %v failed`, resp)
	}
}

func TestNewNewModelResponse(t *testing.T) {
	resp := NewNewModelResponse("Demo", "/Demo")
	b, _ := json.Marshal(&resp)
	msg := `{"msg":{"id":"Demo","path":"/Demo"},"jsonType":"NewModel","statusCode":200}`
	if !assert.EqualValues(t, msg, string(b)) {
		t.Errorf(`Deserialisation of %v failed`, resp)
	}
}

func TestNewErrorResponse(t *testing.T) {
	resp := NewErrorResponse("Problem while deserializing")
	b, _ := json.Marshal(&resp)
	msg := `{"msg":"Problem while deserializing","jsonType":"Error","statusCode":500}`
	if !assert.EqualValues(t, msg, string(b)) {
		t.Errorf(`Deserialisation of %v failed`, resp)
	}
}
