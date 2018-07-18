package result

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/brada954/restshell/shell"
)

// AddResultCommand -- Command structure with options
type AddResultCommand struct {
	// Place getopt option value pointers here
	optionLoadXml  *bool
	optionLoadJson *bool
	optionLoadText *bool
}

func NewAddResultCommand() *AddResultCommand {
	return &AddResultCommand{}
}

func (cmd *AddResultCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("file")

	cmd.optionLoadXml = set.BoolLong("xml", 0, "Load as XML content")
	cmd.optionLoadJson = set.BoolLong("json", 0, "Load as JSON content")
	cmd.optionLoadText = set.BoolLong("text", 0, "Load as text content")

	// Add command helpers for verbose, debug, restclient and output formatting
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

// Execute -- Addresult command to load file data like a REST response
func (cmd *AddResultCommand) Execute(args []string) error {
	// Validate arguments
	if len(args) != 1 {
		return shell.ErrArguments
	}

	contentType := "text/plain"
	if *cmd.optionLoadXml {
		contentType = "application/xml"
	} else if *cmd.optionLoadJson {
		contentType = "application/json"
	}

	// Execute commands
	file, err := os.Open(args[0])
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(file)
	data := string(b)

	shell.PushText(contentType, data, err)

	if shell.IsCmdDebugEnabled() || shell.IsCmdVerboseEnabled() {
		if shell.IsStringBinary(data) {
			fmt.Fprintf(shell.OutputWriter(), "Read %d bytes (binary data)\n", len(data))
		} else {
			fmt.Fprintf(shell.OutputWriter(), data)
		}
	} else {
		fmt.Fprintf(shell.OutputWriter(), "Read %d bytes\n", len(data))
		return nil
	}
	return nil
}
