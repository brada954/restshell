package example

import "github.com/brada954/restshell/shell"

func init() {
	AddCommands()
}

func AddCommands() {
	shell.AddCommand("example", shell.CategoryUtilities, NewExampleCommand())
	shell.AddCommand("exquery", shell.CategoryUtilities, NewExqueryCommand())
}
