package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type VerboseCommand struct {
	// Place getopt option value pointers here
}

func NewVerboseCommand() *VerboseCommand {
	return &VerboseCommand{}
}

func (cmd *VerboseCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("value")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *VerboseCommand) Execute(args []string) error {
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "off" {
			shell.SetVerbose(false)
		} else if strings.ToLower(args[0]) == "on" {
			shell.SetVerbose(true)
		} else {
			return errors.New("Invalid value for setting verbose: " + args[0])
		}
	} else {
		verbose := "OFF"
		if shell.IsCmdDebugEnabled() {
			verbose = "ON"
		}
		fmt.Fprintln(shell.OutputWriter(), "Verbose is", verbose)
	}
	return nil
}
