package shell

import (
	"testing"
)

func TestRunInterfaceDefinition(t *testing.T) {
	var cmd Command = NewRunCommand()

	if _, ok := cmd.(Trackable); !ok {
		t.Errorf("Run command was not implementing Trackable interface")
	}
}
