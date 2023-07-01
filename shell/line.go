package shell

import (
	"errors"
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
func NewCommandLine(input string, shellPrefix string) (line *Line, reterr error) {
	line = &Line{}
	reterr = nil

	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				reterr = errors.New(msg)
			} else {
				reterr = fmt.Errorf("unknown parsing error: %v", r)
			}
		}
	}()

	line.OriginalLine = input
	line.Echo = false
	line.Step = false
	line.NoSubstitute = false
	line.IsComment = false
	line.Command = ""
	line.ArgString = ""

	line.CmdLine = strings.TrimSpace(input)

	if len(line.CmdLine) == 0 {
		return
	}

	if strings.HasPrefix(line.CmdLine, "#") {
		line.IsComment = true
		return
	}

	line.handleSpecialCharacters()

	if len(shellPrefix) > 0 {
		text := strings.ToLower(line.CmdLine)
		if text != "q" && text != "quit" && text != "shell" {
			line.CmdLine = strings.TrimSpace(shellPrefix + line.CmdLine)
		}
	}

	line.CmdLine, reterr = ExpandAlias(line.CmdLine)

	if !line.NoSubstitute {
		line.CmdLine = PerformVariableSubstitution(line.CmdLine)
	}

	// Re-calculate new command after aliases and substitutions
	line.handleSpecialCharacters()
	line.splitCommandAndArgs()
	return
}

// GetCmdAndArguments -- get the tokens of the commmand line fully parsed (obsolete)
func (line *Line) GetCmdAndArguments() []string {
	return line.GetTokens()
}

// GetTokens -- get the line tokens
func (line *Line) GetTokens() []string {
	args := LineParse(line.ArgString)
	return append([]string{line.Command}, args...)
}

// splitCommandAndArgs -- Split CmdLine to Command and ArgString
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

// handleSpecialCharacters -- Parse leading special characters
func (line *Line) handleSpecialCharacters() {
	for notDone := true; notDone; {
		if strings.HasPrefix(line.CmdLine, "@") {
			line.Echo = true
			line.CmdLine = strings.TrimSpace(strings.TrimLeft(line.CmdLine, "@"))
		} else if strings.HasPrefix(line.CmdLine, "$") {
			line.NoSubstitute = true
			line.CmdLine = strings.TrimSpace(strings.TrimLeft(line.CmdLine, "$"))
		} else if strings.HasPrefix(line.CmdLine, "!") { // Legacy character to be deprecated
			line.NoSubstitute = true
			line.CmdLine = strings.TrimSpace(strings.TrimLeft(line.CmdLine, "!"))
		} else {
			notDone = false
		}
	}
}
