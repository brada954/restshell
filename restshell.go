package main

import (
	"os"

	"github.com/brada954/restshell/shell"
)

func main() {
	exitCode := shell.RunShell(shell.GetDefaultStartupOptions())
	os.Exit(exitCode)
}
