package result

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/brada954/restshell/shell"
)

// LoadCommand -- Command structure with options
type LoadCommand struct {
	// Place getopt option value pointers here
	optionLoadXml    *bool
	optionLoadJson   *bool
	optionLoadText   *bool
	optionSubstitute *bool
	optionLoadVar    *bool
}

func NewLoadCommand() *LoadCommand {
	return &LoadCommand{}
}

func (cmd *LoadCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("file|variable")

	cmd.optionLoadXml = set.BoolLong("xml", 0, "Load as XML content")
	cmd.optionLoadJson = set.BoolLong("json", 0, "Load as JSON content")
	cmd.optionLoadText = set.BoolLong("text", 0, "Load as text content")
	cmd.optionLoadVar = set.BoolLong("var", 0, "Load the variable value")
	cmd.optionSubstitute = set.BoolLong("subst", 0, "Perform string substitution on loaded content")

	// Add command helpers for verbose, debug, restclient and output formatting
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

// Execute -- Addresult command to load file data like a REST response
func (cmd *LoadCommand) Execute(args []string) error {
	// Validate arguments
	if len(args) != 1 {
		return shell.ErrArguments
	}

	var data string

	// Execute commands
	if *cmd.optionLoadVar {
		data = shell.GetGlobalString(args[0])
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		data = string(b)
	}

	if *cmd.optionSubstitute {
		data = shell.PerformVariableSubstitution(data)
	}

	contentType := "text/plain"
	if *cmd.optionLoadXml {
		contentType = "application/xml"
	} else if *cmd.optionLoadJson {
		contentType = "application/json"
	} else if *cmd.optionLoadText {
		contentType = "text/plain"
	} else {
		// Very rudimentary tests for json and xml (TODO: Expand on)
		d := strings.TrimSpace(data)
		if len(d) >= 3 && d[0] == '<' && d[len(d)-1] == '>' {
			contentType = "application/xml"
		} else if len(d) >= 2 &&
			((d[0] == '{' && d[len(d)-1] == '}') ||
				(d[0] == '[' && d[len(d)-1] == ']')) {
			contentType = "application/json"
		}
	}

	shell.PushText(contentType, data, nil)

	if shell.IsCmdDebugEnabled() || shell.IsCmdVerboseEnabled() {
		if shell.IsStringBinary(data) {
			fmt.Fprintf(shell.OutputWriter(), "Read %d bytes (binary data)\n", len(data))
		} else {
			fmt.Fprintln(shell.OutputWriter(), data)
		}
	} else {
		fmt.Fprintf(shell.OutputWriter(), "Read %d bytes\n", len(data))
		return nil
	}
	return nil
}
