package example

import (
	"fmt"
	"strconv"

	"github.com/brada954/restshell/shell"
)

type BmCommand struct {
	// Place getopt option value pointers here
}

func NewBmCommand() *BmCommand {
	return &BmCommand{}
}

func (cmd *BmCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("value")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdBenchmarks)
}

func (cmd *BmCommand) Execute(args []string) error {
	// Validate arguments

	// Execute commands
	bm := shell.NewBenchmark(shell.GetCmdIterationValue())
	for i, _ := range bm.Iterations {
		bm.StartIteration(i)
		fmt.Fprintln(shell.OutputWriter(), "Iteraction: ", strconv.Itoa(i))
		bm.EndIteration(i)
		bm.SetIterationStatus(i, nil)
	}

	bm.Dump("Example", shell.GetStdOptions(), shell.IsCmdVerboseEnabled())
	return nil
}
