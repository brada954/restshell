package shell

// getopt.go - Provide an interface to vendored getopt to help reduce dependencies on
// external getopt package that is in venor

import (
	"errors"
	"io"

	"github.com/pborman/getopt/v2"
)

// CmdSet -- Interface exposing the supported interfaces to commands
// for setting options
type CmdSet interface {
	Reset()
	Usage()
	SetProgram(string)
	SetParameters(string)
	SetUsage(func())
	PrintUsage(io.Writer)
	BoolLong(string, rune, ...string) *bool
	StringLong(string, rune, string, ...string) *string
	IntLong(string, rune, int, ...string) *int
	Int64Long(string, rune, int64, ...string) *int64
	Args() []string
	Arg(int) string
	NArgs() int
}

func NewCmdSet() CmdSet {
	return getopt.New()
}

// CmdParse -- implements the parse/getopt function by hiding an Option type
// not to be exported in the interface
func CmdParse(set CmdSet, tokens []string) error {
	if s, ok := set.(*getopt.Set); ok {
		return s.Getopt(tokens, nil)
	} else {
		return errors.New("Invalid option package used")
	}
}
