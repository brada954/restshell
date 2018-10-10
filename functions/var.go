package functions

import (
	"github.com/brada954/restshell/shell"
)

func init() {
	shell.RegisterSubstitutionHandler(Quote)
}

var Quote = shell.SubstitutionFunction{
	Name:         "quote",
	Group:        "quote",
	FunctionHelp: "Quote a variable",
	Formats: []shell.SubstitutionItemHelp{
		{"var", "Option is a variable name"},
	},
	OptionDescription: "",
	Options:           nil,
	Function:          QuoteSubstitute,
}

// Quote -- quote the input string including escaping inner quotes and escape sequences
func QuoteSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	if cache == nil {
		if format == "var" {
			value = shell.GetGlobalString(option)
		} else {
			value = option
		}
	}
	return QuoteString(value), value
}

// QuoteString -- Quote a string handling '\"' as well
func QuoteString(input string) string {
	result := `"`
	for _, c := range input {
		if c == '"' {
			result = result + `\"`
		} else if c == '\\' {
			result = result + string(`\\`)
		} else {
			result = result + string(c)
		}
	}
	return result + `"`
}
