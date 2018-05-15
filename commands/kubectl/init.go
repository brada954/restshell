package kubectl

import (
	"github.com/brada954/restshell/shell"
)

func init() {
	AddCommands()
}

func AddCommands() {
	shell.AddCommand("kubectl", "Kubernetes", NewKubectlCommand())
}
