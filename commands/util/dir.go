package util

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type DirCommand struct {
	// Place getopt option value pointers here
}

func NewDirCommand() *DirCommand {
	return &DirCommand{}
}

func (cmd *DirCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("value")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *DirCommand) Execute(args []string) error {
	dirArgs := osDirArgs
	for _, a := range args {
		dirArgs = append(dirArgs, a)
	}

	c := exec.Command(osDirCmd, dirArgs...)

	text, err := c.Output()
	if err != nil {
		return errors.New("Dir Error: " + err.Error())
	}
	fmt.Fprintf(shell.OutputWriter(), "%s\n", string(text))
	return nil
}
