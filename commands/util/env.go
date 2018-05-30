package util

import (
	"fmt"
	"os"

	"github.com/brada954/restshell/shell"
)

type EnvCommand struct {
	// Place getopt option value pointers here
}

func NewEnvCommand() *EnvCommand {
	return &EnvCommand{}
}

func (cmd *EnvCommand) AddOptions(set shell.CmdSet) {
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
