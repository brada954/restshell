package shell

import (
	"reflect"
	"strings"
)

// Command - interface for basic command
type Command interface {
	Execute([]string) error
	AddOptions(CmdSet)
}

// Abortable - interface for commands that support abort
type Abortable interface {
	Abort()
}

// Trackable - interface that overrides tracking mechanisms
type Trackable interface {
	// DoNotCount - prevents the executed command to be counted as a command executed
	DoNotCount() bool
	// DoNotClearError - prevents the command from affecting the LastError state of a previous command
	DoNotClearError() bool
	// CommandCount - returns the number of commands this command may have executed
	CommandCount() int
}

// LineProcessor - interface for commands that execute whole line
type LineProcessor interface {
	ExecuteLine(line string, echoed bool) error
}

// CommandWithSubcommands - interface for commands that have sub-commands
type CommandWithSubcommands interface {
	GetSubCommands() []string
}

// Variables for supported categorizations of commands
var (
	CategoryHttp        = "Http"
	CategorySpecialized = "Specialized"
	CategoryUtilities   = "Utility"
	CategoryBenchmarks  = "Benchmark"
	CategoryTests       = "Test"
	CategoryAnalysis    = "Result Processing"
	CategoryHelp        = "Help"
)

var cmdMap = make(map[string]Command)
var cmdKeys = make(map[string][]string)
var cmdCategories = make([]string, 0)
var cmdSubCommands = make(map[string][]string)

// AddCommand -- Add a command to registry
// Cmd structures should avoid pointers to data structures so cmd structures can
// be duplicated into separate instances without data collision
func AddCommand(name string, category string, cmd Command) {
	name = strings.ToUpper(name)
	category = strings.ToLower(category)

	validateCmdEntry(name, cmd)
	ensureCategory(category)

	keys, ok := cmdKeys[category]
	if !ok {
		panic("category should exist")
	}

	cmdKeys[category] = append(keys, name)
	cmdMap[name] = cmd

	if subCmd, ok := cmd.(CommandWithSubcommands); ok {
		subcommands := subCmd.GetSubCommands()
		if len(subcommands) > 0 {
			cmdSubCommands[name] = subcommands
		}
	}
}

func ensureCategory(category string) {
	category = strings.ToLower(category)
	if _, ok := cmdKeys[category]; !ok {
		cmdCategories = append(cmdCategories, category)
		cmdKeys[category] = make([]string, 0)
	}
}

func validateCmdEntry(name string, cmd Command) {
	cmdType := reflect.TypeOf(cmd)
	for k, v := range cmdMap {
		if k == name || (v != nil && reflect.TypeOf(v) == cmdType) {
			panic("Command added more than once: " + name)
		}
	}
}

func CommandProcessesLine(cmd interface{}) bool {
	if _, isLine := cmd.(LineProcessor); isLine {
		return true
	}
	return false
}
