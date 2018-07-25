package shell

import (
	"fmt"
	"strings"
)

// Line -- a structure for a parsed line
type Line struct {
	OriginalLine string
	CmdLine      string
	Echo         bool
	Step         bool
	NoSubstitute bool
	IsComment    bool
	Command      string
	ArgString    string
}

// NewCommandLine -- parse command line into a Line
func NewCommandLine(input string, shellPrefix string) (line *Line) {
	line = &Line{}

	line.OriginalLine = input
	line.Echo = false
	line.Step = false
	line.NoSubstitute = false
	line.IsComment = false
	line.Command = ""
	line.ArgString = ""

	line.CmdLine = strings.TrimSpace(input)
	if strings.HasPrefix(line.CmdLine, "#") {
		line.IsComment = true
		return
	}

	// Process prefixes, Triming them one at a time
	for notDone := true; notDone; {
		if strings.HasPrefix(line.CmdLine, "@") {
			line.Echo = true
			line.CmdLine = strings.TrimSpace(strings.TrimLeft(line.CmdLine, "@"))
		} else if strings.HasPrefix(line.CmdLine, "!") {
			line.NoSubstitute = true
			line.CmdLine = strings.TrimSpace(strings.TrimLeft(line.CmdLine, "!"))
		} else {
			notDone = false
		}
	}

	if len(shellPrefix) > 0 {
		line.CmdLine = strings.TrimSpace(shellPrefix + line.CmdLine)
	}

	line.splitCommandAndArgs()

	// Perform alias substitution
	if alias, err := GetAlias(line.Command); err == nil {
		if IsDebugEnabled() {
			fmt.Fprintf(ConsoleWriter(), "Using alias: %s\n", alias)
		}
		if len(line.ArgString) > 0 {
			line.CmdLine = alias + " " + line.ArgString
		} else {
			line.CmdLine = alias
		}
	}

	if !line.NoSubstitute {
		line.CmdLine = PerformVariableSubstitution(line.CmdLine)
	}

	// Re-calculate new command after aliases and substitutions
	line.splitCommandAndArgs()
	return
}

// GetCmdAndArguments -- get the tokens of the commmand line fully parsed
func (line *Line) GetCmdAndArguments() []string {
	args := LineParse(line.ArgString)
	return append([]string{line.Command}, args...)
}

// splitCommandAndArgs() -- Regenerate the command and argument
// strings of the line structure based on current CmdLine
func (line *Line) splitCommandAndArgs() {
	line.Command = ""
	line.ArgString = ""

	args := strings.SplitN(line.CmdLine, " ", 2)
	if len(args) > 0 {
		line.Command = strings.ToUpper(args[0])
	}
	if len(args) > 1 {
		line.ArgString = strings.TrimSpace(args[1])
	}
}