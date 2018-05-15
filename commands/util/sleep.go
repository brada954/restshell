package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type SleepCommand struct {
	// Place getopt option value pointers here
	wait chan error
}

func NewSleepCommand() *SleepCommand {
	return &SleepCommand{}
}

func (cmd *SleepCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("[msec]")
	shell.AddCommonCmdOptions(set, shell.CmdVerbose)
}

func (cmd *SleepCommand) Execute(args []string) error {
	// Validate arguments
	var timeArg string = "1000"
	if len(args) > 0 {
		timeArg = args[0]
	}

	value, err := shell.ParseDuration(timeArg, "ms")
	if err != nil {
		return err
	}

	// Execute commands
	if shell.IsCmdVerboseEnabled() {
		if value > 0 {
			fmt.Fprintf(shell.OutputWriter(), "Sleeping for %s...\n", shell.FormatMsTime(float64(value)/float64(time.Millisecond)))
		}
	}

	cmd.wait = make(chan error)
	select {
	case <-cmd.wait:
		if shell.IsCmdDebugEnabled() {
			fmt.Fprintln(shell.ConsoleWriter(), "Sleep aborted")
		}
	case <-time.After(value):
	}
	cmd.wait = nil
	return nil
}

func (cmd *SleepCommand) Abort() {
	c := cmd.wait
	if c != nil {
		c <- errors.New("Aborted")
	}
}
