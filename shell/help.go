package shell

import (
	"fmt"
	"strings"

	"github.com/pborman/getopt/v2"
)

func DisplayHelp() {
	var programName = ProgramName
	var text = `%[1]s [COMMAND [OPTIONS]...]

%[1]s is a command line driven shell to execute commands and tests against
REST APIs. %[1]s can be invoked with arguments respresenting a
single command to execute or without arguments to run in shell mode.

General utility commands like get and post may be used against any REST API
while specialized commands may provide options for interacting with 
specific APIs.

Assertion commands and a script execution engine enable this tool to run 
complicated test scripts.

For more information consult the repository README.md file
`

	fmt.Fprintf(ConsoleWriter(), text, programName)
	for _, category := range cmdCategories {
		fmt.Fprintf(ConsoleWriter(), "\n%s commands:\n", strings.Title(category))
		for _, cmd := range ColumnizeTokens(cmdKeys[category], 5, 12) {
			fmt.Fprintf(ConsoleWriter(), "  %s\n", cmd)
		}
	}

	var finaltext = `Command modifiers when prefixing command:

#  Comment character to ignore the content on the line (must be first)
@  Echo character to display the executing command including 
   expanded variables and aliases
$  Skip variable substitution on the command line
`

	fmt.Fprintln(ConsoleWriter())
	fmt.Fprintln(ConsoleWriter(), finaltext)
}

func DisplayCmdHelp(set *getopt.Set, cmd string) {
	fmt.Fprintf(ConsoleWriter(), "Command Help")
	set.PrintUsage(ConsoleWriter())
}
