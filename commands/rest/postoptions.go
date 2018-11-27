package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/brada954/restshell/shell"
)

const (
	DefaultJsonVar  = ""
	DefaultJsonBody = ""
	DefaultJsonFile = ""
	DefaultFormBody = ""
	DefaultXMLVar   = ""
	DefaultXMLFile  = ""
	DefaultFormVar  = ""
	DefaultBodyFile = ""
)

type PostOptions struct {
	optionUsePut     *bool
	optionUseOption  *bool
	optionJsonVar    *string
	optionJson       *string
	optionJsonFile   *string
	optionXMLFile    *string
	optionXMLVar     *string
	optionForm       *string
	optionFormVar    *string
	optionBodyFile   *string
	optionLastResult *bool
}

type PostBody struct {
	body        string
	contentType string
}

// Body -- the body of a post
func (pb *PostBody) Content() string {
	if pb == nil {
		return ""
	}
	return pb.body
}

// ContentType -- the default content type based on data options
func (pb *PostBody) ContentType() string {
	if pb == nil {
		return ""
	}
	return pb.contentType
}

// AddModifierOptions -- Add options for modifiers
func AddPostOptions(set shell.CmdSet) PostOptions {
	options := PostOptions{}
	options.optionUsePut = set.BoolLong("put", 0, "Use PUT method instead of post")
	options.optionUseOption = set.BoolLong("options", 0, "Use OPTIONS method instead of post")
	options.optionJsonVar = set.StringLong("json-var", 0, DefaultJsonVar, "Use a named variable as body of json request", "name")
	options.optionJson = set.StringLong("json", 0, DefaultJsonBody, "Send the given json in the body", "json")
	options.optionForm = set.StringLong("form", 0, DefaultFormBody, "Send the given form body", "form")
	options.optionFormVar = set.StringLong("form-var", 0, DefaultFormVar, "Use a named variable as body of form", "name")
	options.optionJsonFile = set.StringLong("json-file", 0, DefaultJsonFile, "Use the given file for json request", "file")
	options.optionXMLVar = set.StringLong("xml-var", 0, DefaultXMLVar, "Use a named variable as body of XML request", "name")
	options.optionXMLFile = set.StringLong("xml-file", 0, DefaultXMLFile, "Use the given file for xml request", "file")
	options.optionBodyFile = set.StringLong("body", 0, DefaultBodyFile, "Send the given file in the body", "file")
	options.optionLastResult = set.BoolLong("result", 0, "Use last result in post body")
	return options
}

func getPostBodyFromFile(filename, extension, contentType string) (*PostBody, error) {
	body, err := shell.GetFileContentsOfType(filename, extension)
	if err != nil {
		return nil, err
	}
	return &PostBody{body: body, contentType: contentType}, nil
}

// GetPostBody -- Get a post body based on post options
func (p *PostOptions) GetPostBody() (*PostBody, error) {

	if *p.optionJson != DefaultJsonBody {
		return &PostBody{body: *p.optionJson, contentType: "application/json"}, nil
	} else if *p.optionJsonVar != DefaultJsonVar {
		return &PostBody{body: shell.GetGlobalStringWithFallback(*p.optionJsonVar, ""), contentType: "application/json"}, nil
	} else if *p.optionXMLVar != DefaultXMLVar {
		return &PostBody{body: shell.GetGlobalStringWithFallback(*p.optionXMLVar, ""), contentType: "application/xml"}, nil
	} else if *p.optionForm != DefaultFormBody {
		return &PostBody{body: *p.optionForm, contentType: "application/x-www-form-urlencoded"}, nil
	} else if *p.optionJsonFile != DefaultJsonFile {
		return getPostBodyFromFile(*p.optionJsonFile, "json", "application/json")
	} else if *p.optionXMLFile != DefaultXMLFile {
		return getPostBodyFromFile(*p.optionXMLFile, "xml", "application/xml")
	} else if *p.optionBodyFile != DefaultBodyFile {
		return getPostBodyFromFile(*p.optionBodyFile, "txt", "text/plain")
	} else if *p.optionFormVar != DefaultFormVar {
		return &PostBody{body: shell.GetGlobalStringWithFallback(*p.optionFormVar, ""), contentType: "application/x-www-form-urlencoded"}, nil
	} else if *p.optionLastResult {
		r, err := shell.PeekResult(0)
		if err != nil {
			return nil, err
		}
		return &PostBody{body: r.Text, contentType: strings.ToLower(r.ContentType)}, nil
	}
	return nil, errors.New("No post body provided")
}

// GetPostMethod -- Returns the configured HTTP method
func (p *PostOptions) GetPostMethod() string {
	method := http.MethodPost
	if *p.optionUsePut {
		method = http.MethodPut
	} else if *p.optionUseOption {
		method = http.MethodOptions
	}
	return method
}
