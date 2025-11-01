package server

import "testing"

func TestNoop(t *testing.T) {
	result := Noop()
	if !result {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestStart(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Start() panicked: %v", r)
		}
	}()

	Start()
}
