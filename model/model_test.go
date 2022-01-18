package model

import (
	"encoding/json"
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
