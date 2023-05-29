package about

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/brada954/restshell/shell"
)

type SubstitutionTopic struct {
	Key         string
	Title       string
	Description string
	About       string
}

var localSubstitutionTopic = &SubstitutionTopic{
	Key:         "SUBST",
	Title:       "Substitution",
	Description: "Substitution functions in commands and input files",
	About: `Variable substitution is available with the command line and conditionally with commands 
reading input files (e.g. post, load, etc). The SET command as well as other commands may
add variables the global name space. Substitution occurs before the command line is
parsed into parameters.

Use the following syntax when referencing a variable from the command line:

	%%%%varname%%%%

In addition to variable substitution, functions can provide complex substitions
using calculated data or special formatting of variables. Restshell includes some basic
functions and enables developers to add their own functions.

The substitution function syntax rquires following format:

	%%%%funcname(instancekey,format,"options")%%%%

Parameters can be omitted as long as the commas are included.
`,
}

// NewSubstitutionTopic -- return a topic structure for help about Substitutions
func NewSubstitutionTopic() *SubstitutionTopic {
	return localSubstitutionTopic
}

// GetKey -- return the key to the substitution help topic
func (a *SubstitutionTopic) GetKey() string {
	return a.Key
}

// GetTitle() -- return the title of the substitution about topic
func (a *SubstitutionTopic) GetTitle() string {
	return a.Title
}

// GetDescription -- return the description of substitution about topic
func (a *SubstitutionTopic) GetDescription() string {
	return a.Description
}

// WriteAbout -- write the substitution about topic to the provided writer
func (a *SubstitutionTopic) WriteAbout(o io.Writer) error {
	fmt.Fprintf(o, a.About)
	fmt.Fprintf(o, "\n\nFunctions:\n")
	substitutionFunctionHelpList(o)
	fmt.Fprintf(o, "\nRun \"ABOUT SUBST {funcname}\" to get more details\n\n")
	return nil
}

// WriteSubTopic -- Write the Subtopic about information
func (a *SubstitutionTopic) WriteSubTopic(o io.Writer, fname string) error {
	return substitutionFunctionHelp(o, fname)
}

func substitutionFunctionHelp(o io.Writer, funcName string) error {
	if fn, ok := shell.GetSubstitutionFunction(funcName); ok {
		fmt.Fprintf(o, "%s: %s\n", strings.ToUpper(fn.Name), fn.FunctionHelp)
		if len(fn.Formats) > 0 {
			if len(fn.FormatDescription) > 0 {
				fmt.Fprintf(o, "  %s\n", fn.FormatDescription)
			} else {
				fmt.Fprintf(o, "  Format Specifiers:\n")
			}
			for _, f := range fn.Formats {
				fmt.Fprintf(o, "    %s: %s\n", f.Item, f.Description)
			}
		}
		if len(fn.Options) > 0 {
			if len(fn.OptionDescription) > 0 {
				fmt.Fprintf(o, "  %s\n", fn.OptionDescription)
			} else {
				fmt.Fprintf(o, "  Options Specifiers:\n")
			}
			for _, f := range fn.Options {
				fmt.Fprintf(o, "    %s: %s\n", f.Item, f.Description)
			}
		}
		fmt.Fprintf(o, "\nExample:\n  %%%%%s%%%%\n", generateExample(fn))
		list := shell.SortedGroupSubstitutionFunctionList(fn.Group)
		if len(list) > 1 {
			fmt.Fprintf(o, "\nRelated Functions (grouped to share the key data):\n")
			for _, g := range list {
				if g.Name != fn.Name {
					fmt.Fprintf(o, "  %s\n", g.Name)
				}
			}
		}
		fmt.Fprintln(o)
		return nil
	} else {
		return errors.New("function not defined")
	}
}

func substitutionFunctionHelpList(o io.Writer) {
	arr := shell.SortedSubstitutionFunctionList(true)
	for _, v := range arr {
		fmt.Fprintf(o, "%s: %s\n", strings.ToUpper(v.Name), v.FunctionHelp)
	}
}

func generateExample(fn shell.SubstitutionFunction) string {
	if len(fn.Example) > 0 {
		return fn.Example
	}
	return fmt.Sprintf("%s(keyname,format,\"options\")", strings.ToLower(fn.Name))
}
