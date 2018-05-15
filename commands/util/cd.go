package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type CdCommand struct {
	initialDir string
	// Place getopt option value pointers here
	resetDir *bool
}

func NewCdCommand() *CdCommand {
	cmd := &CdCommand{}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(shell.ErrorWriter(), "WARNING: Unable to get current working directory")
		cmd.initialDir = "."
	} else {
		cmd.initialDir = dir
	}

	return cmd
}

func (cmd *CdCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("value")
	cmd.resetDir = set.BoolLong("reset", 'r', "Reset current working directory to initial startup")
}

func (cmd *CdCommand) Execute(args []string) error {
	if *cmd.resetDir {
		if len(args) > 0 {
			fmt.Fprintln(shell.ErrorWriter(), "Arguments are ignored with -reset option")
		}
		err := os.Chdir(cmd.initialDir)
		if err != nil {
			return errors.New("Unable to reset dir to: " + cmd.initialDir)
		}
		return nil
	}

	switch len(args) {
	case 1:
		err := os.Chdir(args[0])
		if err != nil {
			return errors.New("Unable to change to: " + args[0])
		}
	case 0:
		dir, err := os.Getwd()
		if err == nil {
			fmt.Fprintln(shell.OutputWriter(), dir)
		} else {
			return errors.New("Unable to get current working directory")
		}
	default:
		return errors.New("Invalid number of arguments")
	}
	return nil
}
