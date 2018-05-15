package rest

import (
	"net/http"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
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

func (cmd *BmGetCommand) AddOptions(set *getopt.Set) {
	set.SetParameters("[service route]")
	cmd.optionUseHead = set.BoolLong("head", 0, "Use HTTP HEAD method")
	cmd.optionUseDelete = set.BoolLong("delete", 0, "Use HTTP DELETE method")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdUrl, shell.CmdBasicAuth, shell.CmdQueryParamAuth, shell.CmdRestclient, shell.CmdBenchmarks)
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

	bm := shell.ProcessJob(job, nil, &cmd.aborted)

	if authContext == nil || !authContext.IsAuthed() {
		bm.Note = "Not an authenticated run"
	}

	bm.Dump(method, shell.GetStdOptions(), shell.IsCmdVerboseEnabled())
	return nil
}

func (cmd *BmGetCommand) Abort() {
	cmd.aborted = true
}
