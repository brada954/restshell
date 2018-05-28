package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brada954/restshell/shell"
)

type DebugCommand struct {
	// Place getopt option value pointers here
}

func NewDebugCommand() *DebugCommand {
	return &DebugCommand{}
}

func (cmd *DebugCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("value")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *DebugCommand) Execute(args []string) error {
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "off" {
			shell.SetDebug(false)
		} else if strings.ToLower(args[0]) == "on" {
			shell.SetDebug(true)
		} else {
			return errors.New("Invalid value for setting debug: " + args[0])
		}
	} else {
		debug := "OFF"
		if shell.IsCmdDebugEnabled() {
			debug = "ON"
		}
		fmt.Fprintln(shell.OutputWriter(), "Debug is", debug)
	}
	return nil
}
