package rest

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/brada954/restshell/shell"
)

// Constants for keys of stored globals
var (
	RESTBASEURLKEY  = "RestBaseUrl"
	RESTBASEAUTHKEY = "RestBaseAuth"
)

type BaseCommand struct {
	// Place getopt option value pointers here
	clearOption     *bool
	clearAuthOption *bool
}

func NewBaseCommand() *BaseCommand {
	return &BaseCommand{}
}

func (cmd *BaseCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[baseurl]")
	cmd.clearOption = set.BoolLong("clear", 'c', "Clear the base URL")
	cmd.clearAuthOption = set.BoolLong("clear-auth", 0, "Clear the base auth context")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdUrl, shell.CmdBasicAuth, shell.CmdQueryParamAuth)
}

func (cmd *BaseCommand) Execute(args []string) error {
	// Validate arguments
	if *cmd.clearOption {
		shell.RemoveGlobal(RESTBASEURLKEY)
		shell.SetAuthContext(RESTBASEAUTHKEY, nil)
	}

	if *cmd.clearAuthOption {
		shell.SetAuthContext(RESTBASEAUTHKEY, nil)
	}

	var authContext = shell.GetCmdQueryParamAuthContext(shell.GetCmdBasicAuthContext(nil))
	if authContext != nil {
		shell.SetAuthContext(RESTBASEAUTHKEY, authContext)
	}

	if len(args) == 0 && !(*cmd.clearOption || *cmd.clearAuthOption) {
		fmt.Fprintf(shell.OutputWriter(), "Current Base Url: %s\n", shell.GetGlobalStringWithFallback(RESTBASEURLKEY, "{not set}"))
		if auth, err := shell.GetAuthContext(RESTBASEAUTHKEY); err == nil {
			fmt.Fprintf(shell.OutputWriter(), "Current Auth Type: %v:%s\n", reflect.TypeOf(auth), auth.ToString())
		} else {
			fmt.Fprintf(shell.OutputWriter(), "Current Auth Type: undefined\n")
		}
		return nil
	}

	if len(args) > 0 && len(args[0]) > 0 {
		shell.SetGlobal(RESTBASEURLKEY, args[0])
	}
	return nil
}

func GetBaseAuthContext() shell.Auth {
	auth, err := shell.GetAuthContext(RESTBASEAUTHKEY)
	if err != nil {
		return nil
	} else {
		return auth
	}
}

func GenerateBaseUrl(route string) string {
	result := shell.GetGlobalStringWithFallback(RESTBASEURLKEY, "")
	if len(result) == 0 {
		return ""
	}

	if strings.HasPrefix(route, "/") {
		result = strings.TrimRight(result, "/")
	} else {
		if !strings.HasSuffix(result, "/") {
			result = result + "/"
		}
	}
	return result + route
}
