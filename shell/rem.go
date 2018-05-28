/////////////////////////////////////////////////////////
// This package makes a good starting point
//
// REM is actually executed internally in the processor
// as it does not work on parsed line (for now)
/////////////////////////////////////////////////////////
package shell

import (
	"fmt"
)

type RemCommand struct {
}

func NewRemCommand() *RemCommand {
	return &RemCommand{}
}

func (cmd *RemCommand) AddOptions(set CmdSet) {
}

func (cmd *RemCommand) Execute(args []string) error {
	return nil
}

func (cmd *RemCommand) DoNotCount() bool {
	return true
}

func (cmd *RemCommand) DoNotClearError() bool {
	return true
}

func (cmd *RemCommand) CommandCount() int {
	return 0
}

func (cmd *RemCommand) ExecuteLine(line string, echoed bool) error {
	if !echoed && !IsSilentEnabled() {
		fmt.Fprintln(OutputWriter(), line)
	}
	return nil
}

func (cmd *RemCommand) RequestQuit() bool {
	return false
}

func (cmd *RemCommand) RequestNoStep() bool {
	return true
}
