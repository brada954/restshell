package result

import "github.com/brada954/restshell/shell"

func init() {
	AddCommands()
}

func AddCommands() {
	shell.AddCommand("load", shell.CategoryUtilities, NewLoadCommand())
	shell.AddCommand("dump", shell.CategoryUtilities, NewDumpCommand())
}
