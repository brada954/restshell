package example

import (
	"fmt"
	"io"

	"github.com/brada954/restshell/shell"
)

type ExSubCmdCommand struct {
	// options
}

func NewExSubCmdCommand() *ExSubCmdCommand {
	return &ExSubCmdCommand{}
}

func (cmd *ExSubCmdCommand) GetSubCommands() []string {
	var commands = []string{"FIRST", "SECOND"}
	return shell.SortedStringSlice(commands)
}

func (cmd *ExSubCmdCommand) AddOptions(set shell.CmdSet) {
	set.SetProgram("exsubcmd [sub command]")
	set.SetUsage(func() {
		cmd.HeaderUsage(shell.ConsoleWriter())
		set.PrintUsage(shell.ConsoleWriter())
		cmd.ExtendedUsage(shell.ConsoleWriter())
	})
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *ExSubCmdCommand) HeaderUsage(w io.Writer) {
	fmt.Fprintln(w, "EXSUBCMD SUBCMD")
	fmt.Fprintln(w)
	fmt.Fprintln(w, `Example command to demostrate sub commads`)
	fmt.Fprintln(w)
}

func (cmd *ExSubCmdCommand) ExtendedUsage(w io.Writer) {
	fmt.Fprintf(w, "\nSub Commands\n")
	lines := shell.ColumnizeTokens(cmd.GetSubCommands(), 4, 15)
	for _, v := range lines {
		fmt.Fprintf(w, "  %s\n", v)
	}
}

func (cmd *ExSubCmdCommand) Execute(args []string) error {

	if len(args) != 1 {
		return shell.ErrArguments
	}

	if !shell.ContainsCommand(args[0], cmd.GetSubCommands()) {
		return shell.ErrArguments
	}

	switch args[0] {
	case "FIRST":
		fmt.Fprintln(shell.ConsoleWriter(), "First command")
	case "SECOND":
		fmt.Fprintln(shell.ConsoleWriter(), "Second Command")
	default:
		return shell.ErrInvalidSubCommand
	}
	return nil
}
