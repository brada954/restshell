package rest

import (
	"net/http"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

type GetCommand struct {
	// Place getopt option value pointers here
	optionUseHead   *bool
	optionUseDelete *bool
}

func NewGetCommand() *GetCommand {
	return &GetCommand{}
}

func (cmd *GetCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("[service route]")
	cmd.optionUseHead = set.BoolLong("head", 0, "Use HTTP HEAD method")
	cmd.optionUseDelete = set.BoolLong("delete", 0, "Use HTTP DELETE method")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdSilent, shell.CmdUrl, shell.CmdBasicAuth, shell.CmdQueryParamAuth, shell.CmdRestclient)
}

func (cmd *GetCommand) Execute(args []string) error {
	// Determine route
	route := ""
	if len(args) > 0 {
		route = args[0]
	}

	// Build URL
	url := shell.GetCmdUrlValue(GenerateBaseUrl(route))
	if url == "" {
		return shell.PushError(shell.ErrArguments)
	}

	method := http.MethodGet
	if *cmd.optionUseHead {
		method = http.MethodHead
	} else if *cmd.optionUseDelete {
		method = http.MethodDelete
	}

	// Get an auth context
	var authContext = shell.GetCmdBasicAuthContext(shell.GetCmdQueryParamAuthContext(GetBaseAuthContext()))

	// Execute commands
	client := shell.NewRestClientFromOptions()
	result, err := client.DoMethod(method, authContext, url)
	return shell.RestCompletionHandler(result, err, nil)
}
