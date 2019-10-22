package rest

import (
	"net/http"

	"github.com/brada954/restshell/shell"
)

type BmGetCommand struct {
	// Place getopt option value pointers here
	optionUseHead   *bool
	optionUseDelete *bool
	// Processing variables
	aborted bool
}

func NewBmGetCommand() *BmGetCommand {
	return &BmGetCommand{}
}

func (cmd *BmGetCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[service route]")
	cmd.optionUseHead = set.BoolLong("head", 0, "Use HTTP HEAD method")
	cmd.optionUseDelete = set.BoolLong("delete", 0, "Use HTTP DELETE method")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdUrl, shell.CmdBasicAuth,
		shell.CmdQueryParamAuth, shell.CmdRestclient, shell.CmdBenchmarks, shell.CmdTimeout)
}

func (cmd *BmGetCommand) Execute(args []string) error {
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

	method := http.MethodGet
	if *cmd.optionUseHead {
		method = http.MethodHead
	} else if *cmd.optionUseDelete {
		method = http.MethodDelete
	}

	// Get an auth context
	var authContext = shell.GetCmdBasicAuthContext(shell.GetCmdQueryParamAuthContext(GetBaseAuthContext()))

	// Execute command using the job which supports
	// iterations and concurrency
	var client = shell.NewRestClientFromOptions()

	jobMaker := func() shell.JobFunction {
		rc := &client
		if shell.IsCmdReconnectEnabled() {
			tmprc := shell.NewRestClientFromOptions()
			rc = &tmprc
		}
		return func() (*shell.RestResponse, error) {
			return rc.DoMethod(method, authContext, url)
		}
	}

	o := shell.GetJobOptionsFromParams()
	o.CancelPtr = &cmd.aborted
	o.JobMaker = jobMaker
	if o.Iterations == 0 {
		o.Iterations = 10
	}

	bm := shell.NewBenchmark(o.Iterations)
	if authContext == nil || !authContext.IsAuthed() {
		bm.Note = "Not an authenticated run"
	}

	shell.ProcessJob(o, bm)
	bm.Dump(method, shell.GetStdOptions(), shell.IsCmdVerboseEnabled())
	return nil
}

func (cmd *BmGetCommand) Abort() {
	cmd.aborted = true
}
