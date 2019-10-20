package functions

import (
	"strconv"
	"strings"

	"github.com/brada954/restshell/shell"
)

func init() {
	shell.RegisterSubstitutionHandler(NewIteratorDefinition)
}

// NewIteratorDefinition -- definition of iterator substitute function
var NewIteratorDefinition = shell.SubstitutionFunction{
	Name:         "iterator",
	Group:        "iterator",
	FunctionHelp: "Return current value of a named (variable) iterator",
	Formats: []shell.SubstitutionItemHelp{
		{"varname", "Variable to be iterated"},
	},
	OptionDescription: "",
	Options: []shell.SubstitutionItemHelp{
		{"{increment}", "(Default = 1) Provide the integer increment"},
	},
	Function: NewIteratorSubstitute,
}

// NewIteratorSubstitute -- Implementatino of iterator substitution
func NewIteratorSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	var iterator int64
	var increment int64 = 1
	var key = format
	var options = strings.Split(option, ",")

	if len(options) > 0 && len(options[0]) > 0 {
		if incrParam, err := strconv.ParseInt(options[0], 10, 32); err != nil {
			increment = 1
			panic("Invalid increment in function: " + subname + "  Var: " + format)
		} else {
			increment = incrParam
		}
	}

	if cache == nil {
		iteratorString := shell.GetGlobalStringWithFallback(key, "0")
		if v, err := strconv.ParseInt(iteratorString, 10, 64); err != nil {
			iterator = 0
		} else {
			iterator = v
			shell.SetGlobal(key, strconv.FormatInt(iterator+increment, 10))
		}

	} else {
		iterator = cache.(int64)
	}

	switch format {
	default:
		return strconv.FormatInt(iterator, 10), iterator
	}
}
