package shell

import "testing"

func TestCreateInterface(t *testing.T) {
	i := NewCmdOptions()

	i.Reset()

	i.GetOpt(args []string, fn func(Options) bool)
}
