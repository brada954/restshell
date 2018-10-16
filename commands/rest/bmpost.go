package rest

import (
	"strings"

	"github.com/brada954/restshell/shell"
)

type BmPostCommand struct {
	// Place getopt option value pointers here
	useSubstitution             *bool
	useSubstitutionPerIteration *bool
	postOptions                 PostOptions
	// Processing variables
	aborted bool
}

func NewBmPostCommand() *BmPostCommand {
	return &BmPostCommand{}
}

func (cmd *BmPostCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[service route]")
	cmd.useSubstitution = set.BoolLong("subst", 0, "Run variable substitution on initial post data")
	cmd.useSubstitutionPerIteration = set.BoolLong("subst-per-call", 0, "Run variable substitution on post data for each post")
	cmd.postOptions = AddPostOptions(set)
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdUrl, shell.CmdBasicAuth, shell.CmdQueryParamAuth, shell.CmdRestclient, shell.CmdBenchmarks)
}

func (cmd *BmPostCommand) Execute(args []string) error {
	cmd.aborted = false

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

	method := cmd.postOptions.GetPostMethod()
	postBody, err := cmd.postOptions.GetPostBody()
	if err != nil {
		return err
	}

	body := postBody.Content()
	if *cmd.useSubstitution {
		body = shell.PerformVariableSubstitution(body)
	}

	// Get an auth context
	var authContext = shell.GetCmdBasicAuthContext(shell.GetCmdQueryParamAuthContext(GetBaseAuthContext()))

	// Execute command using the job processor which supports
	// iterations and concurrency
	var client = shell.NewRestClientFromOptions()
	jobMaker := func() shell.JobProcessor {
		rc := &client
		if shell.IsCmdReconnectEnabled() {
			tmprc := shell.NewRestClientFromOptions()
			rc = &tmprc
		}

		postdata := body
		if *cmd.useSubstitutionPerIteration {
			postdata = shell.PerformVariableSubstitution(postdata)
		}

		return func() (*shell.RestResponse, error) {
			if strings.HasSuffix(postBody.ContentType(), "xml") {
				return rc.DoWithXml(method, authContext, url, postdata)
			} else if strings.HasSuffix(postBody.ContentType(), "json") {
				return rc.DoWithJson(method, authContext, url, postdata)
			} else {
				return rc.DoWithForm(method, authContext, url, postdata)
			}
		}
	}

	// Need to consider having a make job function
	// that can do pre-processing

	bm := shell.ProcessJob(jobMaker, nil, &cmd.aborted)

	if authContext == nil || !authContext.IsAuthed() {
		bm.Note = "Not an authenticated run"
	}

	bm.Dump(method, shell.GetStdOptions(), shell.IsCmdVerboseEnabled())
	return nil
}

func (cmd *BmPostCommand) Abort() {
	cmd.aborted = true
}
