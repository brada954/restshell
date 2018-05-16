package util

import (
	"errors"
	"fmt"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type AliasCommand struct {
	// Place getopt option value pointers here
}

func NewAliasCommand() *AliasCommand {
	return &AliasCommand{}
}

func (cmd *AliasCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("[[alias] command]")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *AliasCommand) Execute(args []string) error {
	if len(args) == 0 {
		cmd.displayAliases()
	} else if len(args) == 1 {
		cmd.displayAlias(args[0])
	} else if len(args) == 2 {
		return shell.AddAlias(args[0], args[1], true)
	} else {
		return errors.New("Invalid number of arguments")
	}
	return nil
}

func (cmd *AliasCommand) displayAliases() {
	fmt.Fprintln(shell.OutputWriter(), "Aliases:")
	for _, v := range shell.GetAllAliasKeys() {
		cmd.displayAlias(v)
	}
}

func (cmd *AliasCommand) displayAlias(key string) {
	if alias, err := shell.GetAlias(key); err == nil {
		fmt.Fprintf(shell.OutputWriter(), "%s=%s\n", key, alias)
	} else {
		fmt.Fprintf(shell.ErrorWriter(), "Invalid key for displaying alias: %s\n", key)
	}
}
