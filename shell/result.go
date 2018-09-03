package shell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/subchen/go-xmldom"
)

// Result -- a result that can be placed in the history buffer
// and used by assertion handlers
type Result struct {
	Text             string
	Map              interface{}
	XMLDocument      *xmldom.Document
	Error            error
	HttpStatus       int
	HttpStatusString string
	ContentType      string
	HeaderMap        map[string]string
	CookieMap        map[string]string
	AuthMap          interface{}
	cookies          []*http.Cookie
}

// GetObjectMap -- short cut to get an object if the result was a json object
// returns false if it is not a json(-like) object
func (r *Result) GetObjectMap() (map[string]interface{}, bool) {
	if m, ok := r.Map.(map[string]interface{}); ok {
		return m, true
	}
	return nil, false
}

// GetArrayMap -- short cut to get an array if the result was a json array
// returns false if it is not a json(-like) object
func (r *Result) GetArrayMap() ([]interface{}, bool) {
	if m, ok := r.Map.([]interface{}); ok {
		return m, true
	}
	return nil, false
}

func (r *Result) addCookieMap(resp *RestResponse) error {
	cookies := make(map[string]string, 0)
	r.cookies = resp.GetCookies()
	for _, cookie := range r.cookies {
		cookies[cookie.Name] = cookie.Value
	}
	r.HeaderMap = cookies
	return nil // TODO: Are there any error conditions
}

func (r *Result) addHeaderMap(resp *RestResponse) error {
	headers := make(map[string]string, 0)
	for n, values := range resp.GetHeader() {
		headers[n] = values[0]

		// Special code that evalutes authorization header as a JWT
		// This may need revisiting to be more acceptable of possibiltiies
		if n == "Authorization" {
			authmap, err := decodeJwtClaims(values[0])
			if err == nil {
				if IsCmdDebugEnabled() {
					fmt.Fprintf(ConsoleWriter(), "Pushing AuthMap:\n%v\n", authmap)
				}
				r.AuthMap = authmap
			} else {
				if IsCmdDebugEnabled() {
					fmt.Fprintln(ConsoleWriter(), "Unable to decode JWT")
				}
			}
		}
	}
	r.HeaderMap = headers
	return nil // TODO: Are there any error conditions
}

func (r *Result) addParsedContentToResult(contentType string, data string) {
	r.ContentType = contentType

	switch getResultTypeFromContentType(r.ContentType) {
	case "xml":
		{
			doc, err := makeXMLDOM(data)
			if err != nil {
				fmt.Fprintln(ErrorWriter(), "WARNING: XML ERROR: ", err)
				r.Map = makeRootMap(data)
			} else {
				r.XMLDocument = doc
			}
		}
	case "json":
		{
			resultMap, err := makeResultMapFromJson(data)
			if err != nil {
				fmt.Fprintln(ErrorWriter(), "WARNING: JSON ERROR: ", err)
				resultMap = makeRootMap(data)
			}
			r.Map = resultMap
		}
	case "text":
		r.Map = makeRootMap(data)
	case "html":
		r.Map = makeRootMap(data)
	case "csv":
		r.Map = makeRootMap(data)
	default:
		r.Map = makeRootMap("Unsupported content type returned: " + r.ContentType)
	}
}

func (r *Result) DumpCookies(w io.Writer) {
	for _, v := range r.cookies {
		fmt.Fprintf(w, "Cookie: %s=%s (%v)\n", v.Name, v.Value, v.Expires)
	}
}

func (r *Result) DumpHeader(w io.Writer) {
	for k, v := range r.HeaderMap {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}
}

func (r *Result) DumpResult(w io.Writer, options ...DisplayOption) {
	if IsStatus(options) && !IsHeaders(options) {
		fmt.Fprintf(w, "HEADER: Status(%s)\n", r.HttpStatusString)
	}

	if IsHeaders(options) {
		r.DumpHeader(w)
	}

	if IsCookies(options) {
		r.DumpCookies(w)
	}

	if IsBody(options) {
		if IsStringBinary(r.Text) {
			fmt.Fprintln(w, "Response contains too many unprintable characters to display")
		} else {
			line := r.Text
			if IsPrettyPrint(options) {
				switch getResultTypeFromContentType(r.ContentType) {
				case "xml":
					{
						if doc, err := xmldom.ParseXML(line); err == nil {
							line = doc.XMLPrettyEx("    ")
						}
					}
				case "json":
					{
						var prettyJSON bytes.Buffer
						err := json.Indent(&prettyJSON, []byte(line), "", "\t")
						if err == nil {
							line = prettyJSON.String()
						}
					}
				}
			}
			fmt.Fprintf(w, "Response:\n%s\n", line)
		}
	}
}