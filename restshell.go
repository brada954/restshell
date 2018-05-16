package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

func main() {
	exitCode := 1
	defer func() {
		if r := recover(); r == nil {
			os.Exit(exitCode)
		} else {
			fmt.Fprintln(shell.ConsoleWriter(), "Panic:", r)
			buf := make([]byte, 1<<16)
			length := runtime.Stack(buf, true)
			fmt.Fprintln(shell.ConsoleWriter(), string(buf[:length]))
			os.Exit(100)
		}
	}()

	getopt.Parse()

	if shell.IsDisplayHelpEnabled() {
		shell.DisplayHelp()
		os.Exit(0)
	}

	runInitScripts(shell.IsDebugEnabled())

	if len(getopt.Args()) == 0 {
		cnt, success := shell.CommandProcessor(">> ", os.Stdin, false, false)
		if !success {
			fmt.Println("Did not return success")
		} else {
			exitCode = shell.LastError
			fmt.Printf("Processed %d commands\n", cnt)
		}
	} else {
		runCmdLine(getopt.Args())
		exitCode = shell.LastError
	}
}

func runInitScripts(debug bool) {
	scriptFile := shell.RestShellInitFile
	runInitScript(scriptFile, debug)

	scriptFile = shell.RestShellUserInitFile
	runInitScript(scriptFile, debug)
}

func runInitScript(scriptFile string, debug bool) {
	if _, err := shell.ValidateScriptExists(scriptFile); err != nil {
		scriptFile = filepath.Join(shell.GetExeDirectory(), scriptFile)
	}

	cmdParts := []string{"run -s"}
	if debug {
		cmdParts = append(cmdParts, "-d")
	}
	cmdParts = append(cmdParts, scriptFile)

	cmdStr := strings.Join(cmdParts, " ")
	_, _ = shell.CommandProcessor("", strings.NewReader(cmdStr), false, true)
}

func runCmdLine(args []string) {
	cmdStr := buildCmdLine(args)
	_, _ = shell.CommandProcessor("", strings.NewReader(cmdStr), false, true)
}

func buildCmdLine(args []string) string {
	for i, v := range args {
		if i == 0 {
			args[i] = v
		} else {
			args[i] = quoteString(v)
		}
	}
	return strings.Join(args, " ")
}

func quoteString(str string) string {
	str = strings.Replace(str, "\\", "\\\\", -1)
	return "\"" + strings.Replace(str, "\"", "\\\"", -1) + "\""
}
