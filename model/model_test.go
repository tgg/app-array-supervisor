package model

import (
	"encoding/json"
	"reflect" // Needed for comparison of slices
	"testing"
)

func TestStatus(t *testing.T) {
	var tests = []struct {
		actual   Status
		expected string
	}{
		{StatusUnknown, "Unknown"},
		{StatusStarting, "Starting"},
		{StatusStarted, "Started"},
		{StatusRunning, "Running"},
		{StatusStopping, "Stopping"},
		{StatusStopped, "Stopped"},
	}

	for _, test := range tests {
		if string(test.actual) != test.expected {
			t.Errorf(`%q is not %q`, test.actual, test.expected)
		}
	}
}

func TestTagMapMarshal(t *testing.T) {
	tags := TagMap{
		"lang": "C++",
		"url":  "here",
	}

	b, err := json.Marshal(tags)
	if err != nil {
		t.Errorf(`Marshalling failed: %q`, err)
	}

	if `{"lang":"C++","url":"here"}` != string(b) {
		t.Errorf(`{"lang":"C++","url":"here"} serialized as %s`, string(b))
	}
}

func TestEnvironmentMarshalUnmarshal(t *testing.T) {
	env := Environment{
		Context{
			"this": []string{"that"},
		},
		"This environment",
	}

	b, err := json.Marshal(&env)
	if err != nil {
		t.Errorf(`Marshalling failed: %q`, err)
	}

	if `{"id":["This environment"],"this":["that"]}` != string(b) {
		t.Errorf(`{"id":["This environment"],"this":["that"]} serialized as %s`, string(b))
	}

	var env2 Environment
	err = json.Unmarshal(b, &env2)

	if !reflect.DeepEqual(env, env2) {
		t.Errorf(`%v differs from %v`, env, env2)
	}
}

func TestApplicationUnmarshal(t *testing.T) {
	const m string = `{"id":"FOApp","type":"application","components":[{"id":"Database","type":"component","tags":{"group":"core","type":"database"},"provides":[{"id":"raw data","kind":6}]},{"id":"EventBus","type":"component","tags":{"group":"core"},"commands":{"start":{"type":"javascript","steps":["StartComponent"]},"stop":{"type":"javascript","steps":["StopComponent"]}},"provides":[{"id":"raw events","kind":6}]},{"id":"Cache","type":"component","tags":{"group":"core"},"consumes":["raw events","raw data"]},{"id":"PositionService","type":"component","tags":{"group":"TradePosition"},"provides":[{"id":"/api/Position","object":"Position","kind":2,"protocol":"REST"}],"consumes":["raw events","raw data"]},{"id":"Spreadsheet","type":"component","tags":{"group":"TradePosition"},"consumes":["/api/Position"]}]}`
	var a Application
	err := json.Unmarshal([]byte(m), &a)
	if err != nil {
		t.Errorf(`Deserialisation of %v failed`, m)
	}
}
