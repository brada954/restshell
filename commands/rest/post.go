package rest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brada954/restshell/shell"
)

// PostCommand -- State and options for PostCommand
type PostCommand struct {
	useAuthContext  shell.Auth
	useSubstitution *bool
	postOptions     PostOptions
}

func NewPostCommand() *PostCommand {
	return &PostCommand{}
}

func (cmd *PostCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[service route]")

	cmd.postOptions = AddPostOptions(set)
	cmd.useSubstitution = set.BoolLong("subst", 0, "Perform variable substitution")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdSilent, shell.CmdUrl,
		shell.CmdBasicAuth, shell.CmdRestclient, shell.CmdFormatOutput)
}

// Execute -- Execute the post command
func (cmd *PostCommand) Execute(args []string) error {
	// Determine route
	route := ""
	if len(args) > 0 {
		route = args[0]
	}

	// Build URL
	url := shell.GetCmdUrlValue(GenerateBaseUrl(route))
	if url == "" {
		return shell.PushError(errors.New("unable to construct URL"))
	}

	// Get an auth context
	cmd.useAuthContext = shell.GetCmdBasicAuthContext(shell.GetCmdQueryParamAuthContext(GetBaseAuthContext()))

	method := cmd.postOptions.GetPostMethod()
	postBody, err := cmd.postOptions.GetPostBody()
	if err != nil {
		return err
	}

	body := postBody.Content()
	if *cmd.useSubstitution {
		body = shell.PerformVariableSubstitution(body)
	}

	if shell.IsVariableSubstitutionComplete(body) == false {
		fmt.Fprintf(shell.ErrorWriter(), "WARNING: post body contains unsubstituted variables")
	}

	// Execute commands
	client := shell.NewRestClientFromOptions()
	if strings.HasSuffix(postBody.ContentType(), "xml") {
		resp, err := client.DoWithXml(method, cmd.useAuthContext, url, body)
		return shell.RestCompletionHandler(resp, err, nil)
	} else if strings.HasSuffix(postBody.ContentType(), "json") {
		resp, err := client.DoWithJson(method, cmd.useAuthContext, url, body)
		return shell.RestCompletionHandler(resp, err, nil)
	} else {
		resp, err := client.DoWithForm(method, cmd.useAuthContext, url, body)
		return shell.RestCompletionHandler(resp, err, nil)
	}
}
