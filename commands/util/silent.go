package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type SilentCommand struct {
	// Place getopt option value pointers here
}

func NewSilentCommand() *SilentCommand {
	return &SilentCommand{}
}

func (cmd *SilentCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("value")
}

func (cmd *SilentCommand) Execute(args []string) error {
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "off" {
			shell.SetSilent(false)
		} else if strings.ToLower(args[0]) == "on" {
			shell.SetSilent(true)
		} else {
			return errors.New("Invalid value for setting debug: " + args[0])
		}
	} else {
		value := "OFF"
		if shell.IsCmdSilentEnabled() {
			value = "ON"
		}
		fmt.Fprintln(shell.OutputWriter(), "Silent is", value)
	}
	return nil
}
