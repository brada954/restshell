package rest

import (
	"strings"

	"github.com/brada954/restshell/shell"
)

type BmPostCommand struct {
	// Place getopt option value pointers here
	useSubstitution             *bool
	useSubstitutionPerIteration *bool
	optionExpectedStatus        *int
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
	cmd.optionExpectedStatus = set.IntLong("expect-status", 0, 200, "Expected status from post [default=200]")
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

	// Build the job processor that can perform substitution
	// on each iteration if required
	var client = shell.NewRestClientFromOptions()
	jobMaker := func() shell.JobFunction {
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

	// Execute command using the job processor which supports
	// iterations and concurrency
	o := shell.GetJobOptionsFromParams()
	o.CancelPtr = &cmd.aborted
	o.JobMaker = jobMaker
	o.CompletionHandler = shell.MakeJobCompletionForExpectedStatus(*cmd.optionExpectedStatus)

	bm := shell.NewBenchmark(o.Iterations)
	shell.ProcessJob(o, &bm)

	if authContext == nil || !authContext.IsAuthed() {
		bm.Note = "Not an authenticated run"
	}

	bm.Dump(method, shell.GetStdOptions(), shell.IsCmdVerboseEnabled())
	return nil
}

func (cmd *BmPostCommand) Abort() {
	cmd.aborted = true
}
