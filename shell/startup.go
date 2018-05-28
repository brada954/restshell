package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pborman/getopt/v2"
)

// Default settings for startup
var (
	DefaultInitFileName  = ".rsconfig"
	DefaultInitFileExt   = ".user"
	DefaultScriptFileExt = ".rshell"
	ProgramName          = "RestShell"
	ProgramArgs          = make([]string, 0, 0)
)

// StartupOptions -- configuration available to the shell
type StartupOptions struct {
	DebugInit         bool
	InitFileName      string
	InitFileExt       string
	ScriptFileExt     string
	AbortOnExceptions bool
}

// GetDefaultStartupOptions return an interface to the options for the shell startup
func GetDefaultStartupOptions() StartupOptions {
	return StartupOptions{
		DebugInit:         false,
		InitFileName:      DefaultInitFileName,
		InitFileExt:       DefaultInitFileExt,
		ScriptFileExt:     DefaultScriptExtension,
		AbortOnExceptions: false,
	}
}

// RunShell -- process command line and init scripts
// and run command processor
func RunShell(options StartupOptions) (exitCode int) {
	exitCode = 1
	if options.AbortOnExceptions == false {
		defer func() {
			if r := recover(); r == nil {
				return // Pass-thru existing error code
			} else {
				fmt.Fprintln(ConsoleWriter(), "Panic:", r)
				buf := make([]byte, 1<<16)
				length := runtime.Stack(buf, true)
				fmt.Fprintln(ConsoleWriter(), string(buf[:length]))
				exitCode = 100 // Return 100 for exception
			}
		}()
	}

	if len(options.InitFileName) > 0 {
		DefaultInitFileName = options.InitFileName
	}
	if len(options.InitFileExt) > 0 {
		DefaultInitFileExt = options.InitFileExt
	}
	if len(options.ScriptFileExt) > 0 {
		DefaultScriptFileExt = options.ScriptFileExt
	}

	getopt.Parse()

	if len(os.Args) > 0 {
		ProgramName = os.Args[0]
	}
	ProgramArgs = getopt.Args()

	if IsDisplayHelpEnabled() {
		DisplayHelp()
		return 0
	}

	runInitScripts(options)

	if len(ProgramArgs) == 0 {
		cnt, success := CommandProcessor(">> ", os.Stdin, false, false)
		if !success {
			fmt.Println("Did not return success")
		} else {
			fmt.Printf("Processed %d commands\n", cnt)
			exitCode = LastError
		}
	} else {
		runCmdLine(ProgramArgs)
		exitCode = LastError
	}
	return exitCode
}

func runInitScripts(options StartupOptions) {
	scriptFile := DefaultInitFileName
	runInitScript(scriptFile, options.DebugInit)

	scriptFile = DefaultInitFileName + DefaultInitFileExt
	runInitScript(scriptFile, options.DebugInit)
}

func runInitScript(scriptFile string, debug bool) {
	if _, err := ValidateScriptExists(scriptFile); err != nil {
		scriptFile = filepath.Join(GetExeDirectory(), scriptFile)
	}

	cmdParts := []string{"run -s"}
	if debug {
		cmdParts = append(cmdParts, "-d")
	}
	cmdParts = append(cmdParts, scriptFile)

	cmdStr := strings.Join(cmdParts, " ")
	_, _ = CommandProcessor("", strings.NewReader(cmdStr), false, true)
}

func runCmdLine(args []string) {
	cmdStr := buildCmdLine(args)
	_, _ = CommandProcessor("", strings.NewReader(cmdStr), false, true)
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
