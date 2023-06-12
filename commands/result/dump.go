package result

import (
	"errors"
	"fmt"
	"io"

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
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdFormatOutput)
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

	dispfunc := displayBytesRead
	if shell.GetCmdOutputFileName() != "" ||
		(shell.IsCmdPrettyPrintEnabled() && !shell.IsCmdOutputShortEnabled()) {
		// Force long form output
		dispfunc = nil
	}
	return shell.OutputResult(result, dispfunc)
}

func displayBytesRead(o io.Writer, result shell.Result) error {
	fmt.Fprintf(o, "Result is %d bytes long\n", len(result.Text))
	return nil
}
