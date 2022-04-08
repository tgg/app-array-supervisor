package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNewModelResponseMsg(t *testing.T) {
	resp := NewNewModelResponse("Demo", []string{"/Demo"})
	b, _ := json.Marshal(&resp.Msg)
	msg := `{"id":"Demo","paths":["/Demo"],"msg":""}`
	if !assert.EqualValues(t, msg, string(b)) {
		t.Errorf(`Deserialisation of %v failed`, resp)
	}
}

func TestNewNewModelResponse(t *testing.T) {
	resp := NewNewModelResponse("Demo", []string{"/Demo"})
	b, _ := json.Marshal(&resp)
	msg := `{"msg":{"id":"Demo","paths":["/Demo"],"msg":""},"type":"NewModel","statusCode":200}`
	if !assert.EqualValues(t, msg, string(b)) {
		t.Errorf(`Deserialisation of %v failed`, resp)
	}
}

func TestNewErrorResponse(t *testing.T) {
	resp := NewErrorResponse("Problem while deserializing")
	b, _ := json.Marshal(&resp)
	msg := `{"msg":"Problem while deserializing","type":"Error","statusCode":500}`
	if !assert.EqualValues(t, msg, string(b)) {
		t.Errorf(`Deserialisation of %v failed`, resp)
	}
}

func TestSendCommandResponse(t *testing.T) {
	info := SendCommandInfo{"start", "start.sh"}
	req := SendCommandRequest{info, "Zipper"}
	resp := NewCommandResponse(0, "ok", req)
	b, _ := json.Marshal(&resp)
	msg := `{"msg":{"componentId":"Zipper","status":0,"commandId":"start","command":"start.sh","result":"ok"},"type":"CommandResponse","statusCode":200}`
	if !assert.EqualValues(t, msg, string(b)) {
		t.Errorf(`Deserialisation of %v failed`, resp)
	}
}
