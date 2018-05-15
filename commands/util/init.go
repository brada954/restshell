package util

import "github.com/brada954/restshell/shell"

func init() {
	AddCommands()
}

func AddCommands() {
	shell.AddCommand("assert", shell.CategoryAnalysis, NewAssertCommand())
	shell.AddCommand("set", shell.CategoryUtilities, NewSetCommand())
	shell.AddCommand("alias", shell.CategoryUtilities, NewAliasCommand())
	shell.AddCommand("debug", shell.CategoryUtilities, NewDebugCommand())
	shell.AddCommand("silent", shell.CategoryUtilities, NewSilentCommand())
	shell.AddCommand("diff", shell.CategoryUtilities, NewDiffCommand())
	shell.AddCommand("verbose", shell.CategoryUtilities, NewVerboseCommand())
	shell.AddCommand("dir", shell.CategoryUtilities, NewDirCommand())
	shell.AddCommand("cd", shell.CategoryUtilities, NewCdCommand())
	shell.AddCommand("log", shell.CategoryUtilities, NewLogCommand())
	shell.AddCommand("env", shell.CategoryUtilities, NewEnvCommand())
	shell.AddCommand("sleep", shell.CategoryUtilities, NewSleepCommand())
	shell.AddCommand("pause", shell.CategoryUtilities, NewPauseCommand())
}
