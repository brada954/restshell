package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/brada954/restshell/shell"
)

// Default values for options
const (
	DefaultJsonVar  = ""
	DefaultJsonBody = ""
	DefaultJsonFile = ""
	DefaultFormBody = ""
	DefaultXmlFile  = ""
	DefaultFormVar  = ""
)

type PostCommand struct {
	useAuthContext shell.Auth
	// Place getopt option value pointers here
	optionUsePut     *bool
	optionUseOption  *bool
	optionJsonVar    *string
	optionJson       *string
	optionJsonFile   *string
	optionXmlFile    *string
	optionForm       *string
	optionFormVar    *string
	optionLastResult *bool
}

func NewPostCommand() *PostCommand {
	return &PostCommand{}
}

func (cmd *PostCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[service route]")
	cmd.optionUsePut = set.BoolLong("put", 0, "Use PUT method instead of post")
	cmd.optionUseOption = set.BoolLong("options", 0, "Use OPTIONS method instead of post")
	cmd.optionJsonVar = set.StringLong("json-var", 0, DefaultJsonVar, "Use a named variable as body of json request", "name")
	cmd.optionJson = set.StringLong("json", 0, DefaultJsonBody, "Send the given json in the body", "json")
	cmd.optionForm = set.StringLong("form", 0, DefaultFormBody, "Send the given form body", "form")
	cmd.optionFormVar = set.StringLong("form-var", 0, DefaultFormVar, "Use a named variable as body of form", "name")
	cmd.optionJsonFile = set.StringLong("json-file", 0, DefaultJsonFile, "Use the given file for json request")
	cmd.optionXmlFile = set.StringLong("xml-file", 0, DefaultXmlFile, "Use the given file for xml request")
	cmd.optionLastResult = set.BoolLong("result", 0, "Use last result in post body")

	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose, shell.CmdSilent, shell.CmdUrl, shell.CmdBasicAuth, shell.CmdRestclient)
}

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

	method := http.MethodPost
	if *cmd.optionUsePut {
		method = http.MethodPut
	} else if *cmd.optionUseOption {
		method = http.MethodOptions
	}

	body := ""
	useJson := false
	useXml := false
	if *cmd.optionJson != DefaultJsonBody {
		body = *cmd.optionJson
		useJson = true
	} else if *cmd.optionJsonVar != DefaultJsonVar {
		body = shell.GetGlobalStringWithFallback(*cmd.optionJsonVar, "")
		useJson = true
	} else if *cmd.optionForm != DefaultFormBody {
		body = *cmd.optionForm
	} else if *cmd.optionJsonFile != DefaultJsonFile {
		filename, err := shell.GetValidatedFileName(*cmd.optionJsonFile, "json")
		if err != nil {
			return err
		}

		body, err = shell.GetFileContents(filename)
		if err != nil {
			return err
		}
		body = shell.PerformVariableSubstitution(body)
		useJson = true
	} else if *cmd.optionXmlFile != DefaultXmlFile {
		filename, err := shell.GetValidatedFileName(*cmd.optionXmlFile, ".xml")
		if err != nil {
			return err
		}
		body, err = shell.GetFileContents(filename)
		if err != nil {
			return err
		}
		body = shell.PerformVariableSubstitution(body)
		useXml = true
	} else if *cmd.optionFormVar != DefaultFormVar {
		body = shell.GetGlobalStringWithFallback(*cmd.optionFormVar, "")
	} else if *cmd.optionLastResult {
		r, err := shell.PeekResult(0)
		if err != nil {
			return err
		}
		body = r.Text
		switch strings.ToLower(r.ContentType) {
		case "application/xml":
			useXml = true
		case "application/json":
			useJson = true
		}
	}

	if shell.IsVariableSubstitutionComplete(body) == false {
		fmt.Fprintf(shell.ErrorWriter(), "WARNING: post body contains unsubstituted variables")
	}

	// Execute commands
	client := shell.NewRestClientFromOptions()
	if useXml {
		resp, err := client.DoWithXml(method, cmd.useAuthContext, url, body)
		return shell.RestCompletionHandler(resp, err, nil)
	} else if useJson {
		resp, err := client.DoWithJson(method, cmd.useAuthContext, url, body)
		return shell.RestCompletionHandler(resp, err, nil)
	} else {
		resp, err := client.DoWithForm(method, cmd.useAuthContext, url, body)
		return shell.RestCompletionHandler(resp, err, nil)
	}
}
