package about

import "github.com/brada954/restshell/shell"

func init() {
	AddCommands()
}

func AddCommands() {
	shell.AddCommand("about", shell.CategoryHelp, NewAboutCommand())
	shell.AddCommand("version", shell.CategoryHelp, NewVersionCommand())
	shell.AddCommand("help", shell.CategoryHelp, nil)
}
