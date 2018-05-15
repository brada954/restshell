package shell

import (
	"testing"
)

func TestRemInterfaceDefinition(t *testing.T) {
	var cmd Command = NewRemCommand()

	if _, ok := cmd.(Trackable); !ok {
		t.Errorf("REM command was not implementing Trackable interface")
	}

	if _, ok := cmd.(LineProcessor); !ok {
		t.Errorf("REM command was not implementing line processor")
	}
}
