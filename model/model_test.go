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
		"lang": []string{"C++", "Java"},
		"url":  []string{"here"},
	}

	b, err := json.Marshal(tags)
	if err != nil {
		t.Errorf(`Marshalling failed: %q`, err)
	}

	if `{"lang":["C++","Java"],"url":["here"]}` != string(b) {
		t.Errorf(`{"lang":["C++","Java"],"url":["here"]} serialized as %s`, string(b))
	}
}

func TestEnvironmentMarshal(t *testing.T) {
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
