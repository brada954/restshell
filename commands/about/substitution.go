package about

import (
	"bytes"

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
	About: `Variable substitution is built into the command line and some functions that
read input files (e.g. post, load, etc). Variables set in the global name space
using the SET command or other commands enable command lines to reference
them. The typical syntanx is:

	%%varname%%

In addition to variable substitution, there are functions that can provide
more dynamic data in substitution or formatting manipulation of variables.
The system comes with functions built in displayed below as well as allowing
developers to build reshell with their own functions.

The substitution function syntax is limited to the following format:

	%%funcname(instancekey, format, "options")%%

Parameters can be omitted as long as the commas are in place. Regex expressions
are used to identify functions so any error in syntax will result in the
function being un-substituted.
`,
}

func NewSubstitutionTopic() *SubstitutionTopic {
	return localSubstitutionTopic
}

func (a *SubstitutionTopic) GetKey() string {
	return a.Key
}

func (a *SubstitutionTopic) GetTitle() string {
	return a.Title
}

func (a *SubstitutionTopic) GetDescription() string {
	return a.Description
}

func (a *SubstitutionTopic) GetAbout() string {
	functionInfo := "\n\nFunctions:\n"

	buf := new(bytes.Buffer)
	shell.SubstitutionFunctionHelp(buf)
	help := buf.String()

	return a.About + functionInfo + help
}
