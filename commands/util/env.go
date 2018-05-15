package util

import (
	"fmt"
	"os"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type EnvCommand struct {
	// Place getopt option value pointers here
}

func NewEnvCommand() *EnvCommand {
	return &EnvCommand{}
}

func (cmd *EnvCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("value")
}

func (cmd *EnvCommand) Execute(args []string) error {
	// Validate arguments

	// Execute commands
	for _, v := range os.Environ() {
		fmt.Fprintf(shell.OutputWriter(), "%s\n", v)
	}

	return nil
}
