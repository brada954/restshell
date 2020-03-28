package rest

import (
	"net/http"
	"time"

	"github.com/brada954/restshell/shell"
)

type SmGetCommand struct {
	// Place getopt option value pointers here
	optionUseHead   *bool
	optionUseDelete *bool
	optionBuckets   *int
	// Processing variables
	aborted bool
}

func NewSmGetCommand() *SmGetCommand {
	return &SmGetCommand{}
}

func (cmd *SmGetCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[service route]")
	cmd.optionUseHead = set.BoolLong("head", 0, "Use HTTP HEAD method")
	cmd.optionUseDelete = set.BoolLong("delete", 0, "Use HTTP DELETE method")
	cmd.optionBuckets = set.IntLong("buckets", 'b', 10, "Time slice buckets for metric collection")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdUrl, shell.CmdBasicAuth, shell.CmdQueryParamAuth, shell.CmdRestclient, shell.CmdBenchmarks)
}

func (cmd *SmGetCommand) Execute(args []string) error {
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

	if *cmd.optionBuckets <= 0 || *cmd.optionBuckets >= 1000000 {
		*cmd.optionBuckets = 10
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
	if o.Duration == 0 {
		o.Duration = 10 * time.Second
	}

	sm := shell.NewSiegemark(o.Duration, *cmd.optionBuckets)
	if authContext == nil || !authContext.IsAuthed() {
		sm.Note = "Not an authenticated run"
	}

	shell.ProcessJob(o, sm)
	sm.Dump(method, shell.GetStdOptions(), shell.IsCmdVerboseEnabled())
	return nil
}

// Abort - Requests the current command to abort
func (cmd *SmGetCommand) Abort() {
	cmd.aborted = true
}
