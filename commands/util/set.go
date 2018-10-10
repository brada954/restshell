package util

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brada954/restshell/shell"
)

// Support for unittesting
var osLookupEnv = os.LookupEnv

type SetCommand struct {
	listOption        *bool
	initOnly          *bool
	valueIsPath       *bool
	valueIsAuthPath   *bool
	valueIsCookiePath *bool
	valueIsHeaderPath *bool
	valueIsVar        *bool
	valueIsFile       *bool
	valueIsEnvVar     *bool
	allowEmpty        *bool
	deleteTempOption  *bool
}

func NewSetCommand() *SetCommand {
	return &SetCommand{}
}

func (cmd *SetCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[key[=[value]]...")
	set.SetUsage(func() {
		set.PrintUsage(shell.ConsoleWriter())
		cmd.ExtendedUsage(shell.ConsoleWriter())
	})
	cmd.listOption = set.BoolLong("list", 'l', "List the globals")
	cmd.initOnly = set.BoolLong("init", 'i', "Inialize if not set already")
	cmd.valueIsPath = set.BoolLong("path", 'p', "Use value as a path into history buffer")
	cmd.valueIsAuthPath = set.BoolLong("path-auth", 0, "Use value as a path into history buffer AuthToken")
	cmd.valueIsCookiePath = set.BoolLong("path-cookie", 0, "Use value as a path into history buffer cookies")
	cmd.valueIsHeaderPath = set.BoolLong("path-header", 0, "Use value as a path into history buffer headers")
	cmd.valueIsVar = set.BoolLong("var", 0, "Use the value as variable name to lookup if exists")
	cmd.valueIsEnvVar = set.BoolLong("env", 0, "Use the value to reference an environment variable if exists")
	cmd.valueIsFile = set.BoolLong("file", 0, "Use the value as a file name to read for value")
	cmd.allowEmpty = set.BoolLong("empty", 0, "Allow an empty string for value")
	cmd.deleteTempOption = set.BoolLong("clear-tmp", 0, "Remove all variables starting with $")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *SetCommand) ExtendedUsage(w io.Writer) {
	fmt.Fprintf(w, "\nAdditional Information:\n")
	fmt.Fprintf(w, "\nVariable substitution has powerful functions for generating substitution data as well.\nRun the command \"ABOUT SUBST\" for more info.\n")
}

func (cmd *SetCommand) Execute(args []string) error {
	argCount := len(args)
	if *cmd.deleteTempOption {
		deleteTemporary()
		return nil
	}

	if *cmd.listOption || argCount == 0 || (argCount == 1 && !strings.ContainsRune(args[0], '=')) {
		var fn func(string, interface{}) bool = nil
		if len(args) >= 1 {
			fn = MakeVariableFilter(args...)
		}
		DisplayGlobals(fn)
		return nil
	}

	if len(args) > 0 {
		for _, v := range args {
			processArg(cmd, v)
		}
	}
	return nil
}

func processArg(cmd *SetCommand, arg string) {
	parts := parseArg(arg)
	if len(parts) == 0 {
		fmt.Fprintln(shell.ErrorWriter(), "Warning: skipping invalid argument: "+arg)
		return
	} else if len(parts) == 1 && !*cmd.allowEmpty {
		shell.RemoveGlobal(parts[0])
	} else {
		var value = ""
		if len(parts) > 1 {
			value = parts[1]
		}

		if *cmd.valueIsFile {
			if len(value) == 0 {
				fmt.Fprintln(shell.ErrorWriter(), "Invalid value for file name")
				return
			}
			v, err := shell.GetFileContents(value)
			if err != nil {
				fmt.Fprintln(shell.ErrorWriter(), err.Error())
				return
			}
			value = v
		}

		if *cmd.valueIsVar {
			if len(value) == 0 {
				fmt.Fprintln(shell.ErrorWriter(), "Invalid value for variable name")
				return
			}
			value = shell.GetGlobalString(value)
			if shell.IsDebugEnabled() {
				fmt.Fprintf(shell.ConsoleWriter(), "Setting %s value to variable: %s\n", parts[0], value)
			}
		}

		if *cmd.valueIsEnvVar {
			if len(value) == 0 {
				fmt.Fprintln(shell.ErrorWriter(), "Invalid value for environment variable")
				return
			}

			if v, ok := osLookupEnv(value); ok {
				value = v
			} else {
				value = ""
			}

			if shell.IsDebugEnabled() {
				fmt.Fprintf(shell.ConsoleWriter(), "Setting %s value to environment variable: %s\n", parts[0], value)
			}
		}

		if *cmd.valueIsPath || *cmd.valueIsAuthPath || *cmd.valueIsCookiePath || *cmd.valueIsHeaderPath {
			if len(value) == 0 {
				fmt.Fprintln(shell.ErrorWriter(), "Invalid value for path name")
				return
			}
			var err error
			if *cmd.valueIsAuthPath {
				value, err = shell.GetValueFromAuthHistory(0, value)
			} else if *cmd.valueIsCookiePath {
				value, err = shell.GetValueFromCookieHistory(0, value)
			} else if *cmd.valueIsHeaderPath {
				value, err = shell.GetValueFromHeaderHistory(0, value)
			} else {
				value, err = shell.GetValueFromHistory(0, value)
			}

			if err != nil {
				fmt.Fprintf(shell.ErrorWriter(), "Warning: value not found, skipping argument: %s\n", arg)
				return
			}

			if shell.IsDebugEnabled() {
				fmt.Fprintf(shell.ConsoleWriter(), "Setting %s value to history value: %s\n", parts[0], value)
			}
		}

		if len(value) > 0 || *cmd.allowEmpty {
			if *cmd.initOnly {
				shell.InitializeGlobal(parts[0], value)
			} else {
				shell.SetGlobal(parts[0], value)
			}
			if shell.IsCmdVerboseEnabled() {
				fmt.Fprintf(shell.OutputWriter(), "%s=%s\n", parts[0], shell.GetGlobal(parts[0]))
			}
		}
	}
}

func parseArg(arg string) []string {
	if strings.ContainsRune(arg, '=') {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) >= 2 && parts[1] == "" {
			return parts[:1]
		}
		return parts
	}
	return []string{}
}

func DisplayGlobals(filter func(string, interface{}) bool) {
	fmt.Fprintln(shell.ConsoleWriter(), "Global values:")
	shell.EnumerateGlobals(displayEntry, filter)
}

func MakeVariableFilter(filterArgs ...string) func(string, interface{}) bool {
	filters := make([]string, 0)
	for _, v := range filterArgs {
		filters = append(filters, strings.ToLower(v))
	}
	return func(key string, _ interface{}) bool {
		key = strings.ToLower(key)
		if len(filters) > 0 {
			for _, filter := range filters {
				if strings.HasPrefix(key, filter) {
					return true
				}
			}
			return false
		}
		return true
	}
}

func displayEntry(k string, v interface{}) {
	switch t := v.(type) {
	case string:
		fmt.Fprintf(shell.ConsoleWriter(), "%s=%s\n", k, t)
	case shell.Auth:
		fmt.Fprintf(shell.ConsoleWriter(), "%s=Auth Context(%t)\n", k, t.IsAuthed())
	default:
		fmt.Fprintf(shell.ConsoleWriter(), "%s={unsupported type}\n", k)
	}
}

func deleteTemporary() {
	filter := MakeVariableFilter("$")
	shell.EnumerateGlobals(removeEntry, filter)
}

func removeEntry(k string, v interface{}) {
	shell.RemoveGlobal(k)
}
