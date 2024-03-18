package about

import (
	"fmt"
	"os"

	"github.com/brada954/restshell/shell"
)

var BuildCommit string
var BuildBranch string
var BuildVersion string

type VersionCommand struct {
	Branch  string
	Commit  string
	Version string
}

func NewVersionCommand() *VersionCommand {
	var branch = "local/unknown"
	var commit = "local"
	var version = "0.0.0"

	if len(BuildVersion) > 0 {
		version = BuildVersion
	}
	if len(BuildBranch) > 0 {
		branch = BuildBranch
	}
	if len(BuildCommit) > 0 {
		commit = BuildCommit
	}
	return &VersionCommand{
		Branch:  branch,
		Version: version,
		Commit:  commit,
	}
}

func (cmd *VersionCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *VersionCommand) Execute(args []string) error {
	executable, err := os.Executable()
	if err != nil {
		executable = "Unknown"
	}

	fmt.Fprintf(shell.OutputWriter(),
		"Version: %s Branch: %s Commit: %s\n",
		cmd.Version,
		cmd.Branch,
		cmd.Commit,
	)
	fmt.Fprintf(shell.OutputWriter(),
		"Executable: %s\n",
		executable,
	)
	return nil
}
