package util

import (
	"errors"
	"fmt"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

var (
	DefaultMessage = "Paused... press \"enter\" to continue."
)

type PauseCommand struct {
	// Place getopt option value pointers here
	messageOption *string
	aborted       bool
}

func NewPauseCommand() *PauseCommand {
	return &PauseCommand{}
}

func (cmd *PauseCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("")
	cmd.messageOption = set.StringLong("message", 'm', DefaultMessage, "Message to display")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *PauseCommand) Execute(args []string) error {
	fmt.Fprintf(shell.OutputWriter(), *cmd.messageOption)
	shell.ReadLine()
	if cmd.aborted {
		return errors.New("Command interrupted")
	}
	return nil
}

func (cmd *PauseCommand) Abort() {
	cmd.aborted = true
}
