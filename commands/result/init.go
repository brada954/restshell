package result

import "github.com/brada954/restshell/shell"

func init() {
	AddCommands()
}

func AddCommands() {
	shell.AddCommand("addresult", shell.CategoryUtilities, NewAddResultCommand())
}
