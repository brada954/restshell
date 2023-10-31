package shell

import (
	"fmt"
)

// Logger interface
type Logger interface {
	LogDebug(string, ...any)
	LogVerbose(string, ...any)
}

type log struct {
	verbose bool
	debug   bool
}

func NewLogger(v bool, d bool) Logger {
	return log{
		verbose: v,
		debug:   d,
	}
}

func (l log) LogDebug(format string, args ...any) {
	if l.debug {
		fmt.Fprintf(ConsoleWriter(), format+"\n", args...)
	}
}

func (l log) LogVerbose(format string, args ...interface{}) {
	if l.verbose {
		fmt.Fprintf(ConsoleWriter(), format+"\n", args...)
	}
}
