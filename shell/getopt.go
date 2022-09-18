package shell

// getopt.go - Provide an interface to vendored getopt to help reduce dependencies on
// external getopt package in vendor
//
// Added extensions to support lists of parameters

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
	StringListLong(string, rune, ...string) *StringList
	Args() []string
	Arg(int) string
	NArgs() int
}

type newCmdSet struct {
	*getopt.Set
	cmdUsage func()
}

// NewCmdSet -- Create a command set for handling options
func NewCmdSet() CmdSet {
	g := getopt.New()
	s := &newCmdSet{g, nil}
	return s
}

// CmdParse -- implements the parse/getopt function by hiding an Option type
// not to be exported in the interface
func CmdParse(set CmdSet, tokens []string) error {
	if s, ok := set.(*newCmdSet); ok {
		return s.Getopt(tokens, nil)
	} else {
		return errors.New("invalid option package used")
	}
}

// StringListLong -- implement a string list option
func (c *newCmdSet) StringListLong(name string, short rune, help ...string) *StringList {
	initial := &StringList{
		Values: make([]string, 0),
	}
	c.StringListVarLong(initial, name, short, help...)
	return initial
}

func (c *newCmdSet) StringListVarLong(p getopt.Value, name string, short rune, helpvalue ...string) getopt.Option {
	return c.FlagLong(p, name, short, helpvalue...)
}

func (c *newCmdSet) SetUsage(usage func()) {
	// Need our own copy of usage as we cannot foce the pborman usage to be called
	c.cmdUsage = usage
	c.Set.SetUsage(usage)
}

func (c *newCmdSet) Usage() {
	if c.cmdUsage != nil {
		c.cmdUsage()
	} else {
		c.defaultUsage()
	}
}

func (c *newCmdSet) defaultUsage() {
	c.Set.PrintUsage(ConsoleWriter())
}

type StringList struct {
	Values []string
}

func (sl *StringList) GetValues() []string {
	return sl.Values
}

func (sl *StringList) Count() int {
	return len(sl.Values)
}

func (sl *StringList) Set(value string, opt getopt.Option) error {
	if len(value) > 0 {
		sl.Values = append(sl.Values, value)
	}
	return nil
}

func (sl *StringList) String() string {
	length := len(sl.Values)
	if length == 0 {
		return ""
	} else {
		return sl.Values[length-1]
	}
}
