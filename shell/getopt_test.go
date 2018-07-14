package shell

import "testing"

func TestCmdSetInterface(t *testing.T) {
	i := NewCmdSet()

	i.Reset()

	i.Args()
}
