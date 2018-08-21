package functions

import (
	"strings"

	"github.com/brada954/restshell/shell"
)

var ToLowerDefinition = shell.SubstitutionFunction{
	Name:              "tolower",
	Group:             "tolower",
	FunctionHelp:      "Lower case an options parameter as identified by format",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          ToLowerSubstitute,
}

var ToUpperDefinition = shell.SubstitutionFunction{
	Name:              "toupper",
	Group:             "toupper",
	FunctionHelp:      "Upper case an options parameter as identified by format",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          ToUpperSubstitute,
}

// ToLowerSubstitute -- returns the ToLower value from options parameter with format
// options to use the option parameter in a variable lookup
func ToLowerSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	if cache == nil {
		if format == "var" {
			value = shell.GetGlobalString(option)
		} else {
			value = option
		}
	}
	return strings.ToLower(value), value
}

// ToUpperSubstitute -- returns the ToUpper value from options parameter with format
// options to use the option parameter in a variable lookup
func ToUpperSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	if cache == nil {
		if format == "var" {
			value = shell.GetGlobalString(option)
		} else {
			value = option
		}
	}
	return strings.ToUpper(value), value
}
