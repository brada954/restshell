package example

import (
	"errors"
	"fmt"
	"io"

	"github.com/brada954/restshell/shell"
)

var (
	EXAMPLE_URL_KEY = "Example_Url"
)

type ExqueryCommand struct {
	// Place getopt option value pointers here
}

func NewExqueryCommand() *ExqueryCommand {
	return &ExqueryCommand{}
}

func (cmd *ExqueryCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("value")

	// Add command helpers for verbose, debug, restclient and output formatting
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdUrl, shell.CmdRestclient, shell.CmdFormatOutput)
}

// Execute -- Exquery command to query a well known test data site and demonstrate basic output
// Use variations of -d and -v and the --out-* options to get example output formats
func (cmd *ExqueryCommand) Execute(args []string) error {
	// Validate arguments

	// Note: Fallbacks are not recommended as they can expose secrets, use local user config file for setting
	// environment variables for EXAMPLE_URL_KEY. Similar for url portions that may include api keys.
	var url = shell.GetGlobalStringWithFallback(EXAMPLE_URL_KEY, "https://jsonplaceholder.typicode.com/users")
	{
		url = shell.GetCmdUrlValue(url)
		if len(url) == 0 {
			return errors.New("No URL specified")
		}

		// Example augmenting a base url with specific routes or data
		if len(args) != 0 {
			url = url + "/" + args[0]
		}
	}

	// Execute commands
	client := shell.NewRestClientFromOptions()
	response, err := client.DoGet(nil, url)

	// Test a short display that displays response length
	// Short display can be used to pretty format the output or condense it
	shortDisplay := func(w io.Writer, r shell.Result) error {
		fmt.Fprintf(w, "Length: %d\n", len(r.Text))
		return nil
	}

	// Display output and return errors
	shell.RestCompletionHandler(response, err, shortDisplay)
	return nil
}
