package example

import (
	"github.com/brada954/restshell/shell"
)

type ExampleCommand struct {
	// Place getopt option value pointers here
}

func NewExampleCommand() *ExampleCommand {
	return &ExampleCommand{}
}

func (cmd *ExampleCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("value")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *ExampleCommand) Execute(args []string) error {
	// Validate arguments

	// Execute commands
	return nil
}
