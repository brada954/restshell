package util

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brada954/restshell/shell"
	"github.com/brada954/restshell/shell/modifiers"
)

// Support for unittesting
var osLookupEnv = os.LookupEnv

type SetCommand struct {
	listOption       *bool
	initOnly         *bool
	valueIsVar       *bool
	valueIsFile      *bool
	valueIsEnvVar    *bool
	allowEmpty       *bool
	deleteTempOption *bool
	modifierOptions  modifiers.ModifierOptions
	historyOptions   shell.HistoryOptions
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
	cmd.valueIsVar = set.BoolLong("var", 0, "Use the value as variable name to lookup if exists")
	cmd.valueIsEnvVar = set.BoolLong("env", 0, "Use the value to reference an environment variable if exists")
	cmd.valueIsFile = set.BoolLong("file", 0, "Use the value as a file name to read for value")
	cmd.allowEmpty = set.BoolLong("empty", 0, "Allow an empty string for value")
	cmd.deleteTempOption = set.BoolLong("clear-tmp", 0, "Remove all variables starting with $")
	_ = set.BoolLong("direct", 0, "Use direct value instead of redirecting to variable or file or history")
	cmd.modifierOptions = modifiers.AddModifierOptions(set)
	cmd.historyOptions = shell.AddHistoryOptions(set, shell.AllPaths)
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
			if len(v) > 2 && strings.HasPrefix(v, "--") {
				*cmd.valueIsVar = false
				*cmd.valueIsFile = false
				*cmd.valueIsEnvVar = false
				cmd.historyOptions.ClearPathOptions()
				switch strings.ToLower(v) {
				case "--var":
					*cmd.valueIsVar = true
				case "--file":
					*cmd.valueIsFile = true
				case "--env":
					*cmd.valueIsEnvVar = true
				case "--path":
					cmd.historyOptions.SetPathOption(shell.ResultPath)
				case "--path-header":
					cmd.historyOptions.SetPathOption(shell.HeaderPath)
				case "--path-auth":
					cmd.historyOptions.SetPathOption(shell.AuthPath)
				case "--path-cookie":
					cmd.historyOptions.SetPathOption(shell.CookiePath)
				case "--direct":
				default:
					return fmt.Errorf("Invalid option in arguments: %s", v)
				}
				continue
			}
			processArg(cmd, v)
		}
	}
	return nil
}

func processArg(cmd *SetCommand, arg string) {
	parts := parseArg(arg)
	if len(parts) == 0 || len(parts[0]) == 0 {
		fmt.Fprintln(shell.ErrorWriter(), "Warning: skipping invalid argument: "+arg)
		return
	} else if len(parts) == 1 && !*cmd.allowEmpty {
		shell.RemoveGlobal(parts[0])
		return
	}

	// Setup exit function to perform final common tasks
	var exitError error
	var variable = parts[0]
	var value = parts[1]
	defer func() {
		if exitError != nil {
			fmt.Fprintf(shell.ErrorWriter(), "Warning: %s; clearing variable: %s\n", exitError, variable)
			if len(variable) > 0 && !*cmd.allowEmpty {
				shell.RemoveGlobal(variable)
				return
			}
			value = ""
		}

		if len(value) > 0 || *cmd.allowEmpty {
			var err error

			if *cmd.initOnly {
				err = shell.InitializeGlobal(variable, value)
			} else {
				err = shell.SetGlobal(variable, value)
			}
			if err != nil {
				fmt.Fprintf(shell.ErrorWriter(), "%s: %s\n", err.Error(), variable)
			} else if shell.IsCmdVerboseEnabled() {
				fmt.Fprintf(shell.OutputWriter(), "%s=%s\n", variable, shell.GetGlobal(variable))
			}
		}
		return
	}()

	if *cmd.valueIsFile {
		if len(value) == 0 {
			exitError = fmt.Errorf("Invalid value for file name")
			return
		}
		v, err := shell.GetFileContents(value)
		if err != nil {
			exitError = err
			return
		}
		value = v
	}

	if *cmd.valueIsVar {
		if len(value) == 0 {
			exitError = fmt.Errorf("Invalid value for variable name")
			return
		}
		v, ok := shell.TryGetGlobalString(value)
		if !ok {
			exitError = fmt.Errorf("Variable does not contain a string value")
			return
		}
		value = v
	}

	if *cmd.valueIsEnvVar {
		if len(value) == 0 {
			exitError = fmt.Errorf("Invalid value for environment variable")
			return
		}

		if v, ok := osLookupEnv(value); ok {
			value = v
		} else {
			exitError = fmt.Errorf("Environment variable not found")
			return
		}
	}

	if cmd.historyOptions.IsHistoryPathOptionEnabled() {
		if len(value) == 0 {
			exitError = fmt.Errorf("Invalid value for path name")
			return
		}

		v, err := cmd.historyOptions.GetValueFromHistory(0, value)
		if err != nil {
			exitError = fmt.Errorf("Path Error: %s", err.Error())
			return
		}
		value = v
	}

	if len(value) > 0 || *cmd.allowEmpty {
		valueModifierFunc := modifiers.ConstructModifier(cmd.modifierOptions)
		if v, err := valueModifierFunc(value); err != nil {
			exitError = fmt.Errorf("Modifier Failure: %s", err.Error())
			return
		} else {
			value = fmt.Sprintf("%v", v)
		}
	}

	return
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
