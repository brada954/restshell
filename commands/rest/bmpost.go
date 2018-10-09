package rest

import (
	"net/http"

	"github.com/brada954/restshell/shell"
)

type BmPostCommand struct {
	// Place getopt option value pointers here
	optionUsePut   *bool
	optionUsePatch *bool
	// Processing variables
	aborted bool
}

func NewBmPostCommand() *BmPostCommand {
	return &BmPostCommand{}
}

func (cmd *BmPostCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[service route]")
	cmd.optionUsePut = set.BoolLong("put", 0, "Use HTTP PUT method")
	cmd.optionUsePatch = set.BoolLong("patch", 0, "Use HTTP PATCH method")
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

	method := http.MethodPost
	if *cmd.optionUsePut {
		method = http.MethodPut
	} else if *cmd.optionUsePatch {
		method = http.MethodPatch
	}

	// Get an auth context
	var authContext = shell.GetCmdBasicAuthContext(shell.GetCmdQueryParamAuthContext(GetBaseAuthContext()))

	// Execute command using the job processor which supports
	// iterations and concurrency
	var client = shell.NewRestClientFromOptions()
	job := func() (*shell.RestResponse, error) {
		rc := &client
		if shell.IsCmdReconnectEnabled() {
			tmprc := shell.NewRestClientFromOptions()
			rc = &tmprc
		}

		return rc.DoMethod(method, authContext, url)
	}

	jobMaker := func() shell.JobProcessor {
		return job
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
