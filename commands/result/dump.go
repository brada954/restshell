package result

import (
	"errors"
	"fmt"

	"github.com/brada954/restshell/shell"
)

// DumpCommand -- Command structure with options
type DumpCommand struct {
}

func NewDumpCommand() *DumpCommand {
	return &DumpCommand{}
}

func (cmd *DumpCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("")

	// Add command helpers for verbose, debug, restclient and output formatting
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

// Execute -- Addresult command to load file data like a REST response
func (cmd *DumpCommand) Execute(args []string) error {
	// Validate arguments
	if len(args) != 0 {
		return shell.ErrArguments
	}

	result, err := shell.PeekResult(0)
	if err != nil {
		return errors.New("No result to dump")
	}
	if shell.IsCmdDebugEnabled() || shell.IsCmdVerboseEnabled() {
		if shell.IsStringBinary(result.Text) {
			fmt.Fprintf(shell.OutputWriter(), "Read %d bytes (binary data)\n", len(result.Text))
		} else {
			fmt.Fprintln(shell.OutputWriter(), result.Text)
		}
	} else {
		fmt.Fprintf(shell.OutputWriter(), "Result is %d bytes long\n", len(result.Text))
		return nil
	}
	return nil
}
