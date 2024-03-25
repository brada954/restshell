package shell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/subchen/go-xmldom"
)

// Result -- a result that can be placed in the history buffer
// and used by assertion handlers
type Result struct {
	Text             string
	Error            error
	HttpStatus       int
	HttpStatusString string
	ContentType      string
	BodyMap          HistoryMap
	HeaderMap        HistoryMap
	CookieMap        HistoryMap
	AuthMap          HistoryMap
	cookies          []*http.Cookie
	headers          map[string]string
}

type ResultPayloadType int

// Path options scenarios for different use cases
const (
	ResultPath     ResultPayloadType = 1
	AuthPath       ResultPayloadType = 2
	CookiePath     ResultPayloadType = 3
	HeaderPath     ResultPayloadType = 4
	AllPaths       ResultPayloadType = 8
	AlternatePaths ResultPayloadType = 9 // All paths but default as default is assumed
)

func NewTextResult(text string) *Result {
	result := &Result{Text: text}
	result.addParsedContentToResult("text/plain", text)
	return result
}

func NewJSONResult(text string) *Result {
	result := &Result{Text: text}
	result.addParsedContentToResult("application/json", text)
	result.Error = nil
	return result
}

func (r *Result) DumpCookies(w io.Writer) {
	fmt.Fprintln(w, "Cookies:")
	for _, v := range r.cookies {
		fmt.Fprintf(w, "%s=%s (%v)\n", v.Name, v.Value, v.Expires)
	}
}

func (r *Result) DumpHeader(w io.Writer) {
	fmt.Fprintln(w, "Headers:")
	for k, v := range r.headers {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}
}

func (r *Result) DumpResult(w io.Writer, options ...DisplayOption) {
	verbose := IsCmdVerboseEnabled()

	if IsStatus(options) || IsHeaders(options) || verbose {
		fmt.Fprintf(w, "HTTP Status: %s\n", r.HttpStatusString)
		verbose = true
	}

	if IsHeaders(options) {
		r.DumpHeader(w)
		verbose = true
	}

	if IsCookies(options) {
		r.DumpCookies(w)
		verbose = true
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
						err := json.Indent(&prettyJSON, []byte(line), "", getJsonPrintIndent())
						if err == nil {
							line = prettyJSON.String()
						}
					}
				}
			}
			fmt.Fprint(w, generateResponseLine(line, verbose))
		}
	}
}

func (r *Result) addCookieMap(resp *RestResponse) error {
	r.cookies = resp.GetCookies()

	m := make(map[string]string, 0)
	for _, cookie := range r.cookies {
		m[cookie.Name] = cookie.Value
	}

	var err error
	r.CookieMap, err = NewSimpleHistoryMap(m)
	return err
}

func (r *Result) addHeaderMap(resp *RestResponse) error {
	m := make(map[string]string, 0)
	for n, values := range resp.GetHeader() {
		m[n] = values[0]

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

	r.headers = m
	r.HeaderMap, _ = NewSimpleHistoryMap(m)
	return nil // TODO: Are there any error conditions
}

func (r *Result) addParsedContentToResult(contentType string, data string) {
	r.ContentType = contentType

	switch getResultTypeFromContentType(r.ContentType) {
	case "xml":
		{
			resultMap, err := NewXmlHistoryMap(data)
			if err != nil {
				fmt.Fprintln(ErrorWriter(), "WARNING: XML ERROR: ", err)
				resultMap, _ = NewTextHistoryMap(data)
			}
			r.BodyMap = resultMap
		}
	case "json":
		{
			resultMap, err := NewJsonHistoryMap(data)
			if err != nil {
				fmt.Fprintln(ErrorWriter(), "WARNING: JSON ERROR: ", err)
				resultMap, _ = NewTextHistoryMap(data)
			}
			r.BodyMap = resultMap
		}
	case "text":
		r.BodyMap, _ = NewTextHistoryMap(data)
	case "html":
		r.BodyMap, _ = NewTextHistoryMap(data)
	case "csv":
		r.BodyMap, _ = NewTextHistoryMap(data)
	case ResultContentBinary:
		r.BodyMap, _ = NewTextHistoryMap("")
	default:
		r.BodyMap, _ = NewTextHistoryMap("Unsupported content type returned: " + r.ContentType)
	}
}

// getResultTypeFromResponse -- Get the result type
//
//	xml, json, text, html, css, csv, media, unknown
func getResultTypeFromContentType(contentType string) ResultContentType {
	// Split off parameters
	parts := strings.Split(contentType, ";")
	if len(parts) == 0 {
		return ResultContentUnknown
	}

	contentType = strings.TrimSpace(strings.ToLower(parts[0]))
	if strings.HasPrefix(contentType, "application/xml") ||
		strings.HasSuffix(contentType, "+xml") {
		return ResultContentXml
	} else if strings.HasPrefix(contentType, "application/json") ||
		strings.HasSuffix(contentType, "+json") {
		return ResultContentJson
	} else if strings.HasPrefix(contentType, "text/plain") ||
		strings.HasPrefix(contentType, "text/calendar") ||
		strings.HasPrefix(contentType, "text/css") {
		return ResultContentText
	} else if strings.HasPrefix(contentType, "text/html") {
		return ResultContentHtml
	} else if strings.HasPrefix(contentType, "application/octet-stream") {
		return ResultContentBinary
	} else if strings.Contains(contentType, "text/csv") {
		return ResultContentCsv
	}

	return ResultContentUnknown
}

func generateResponseLine(line string, verbose bool) string {
	responseLabel := ""
	if verbose {
		responseLabel = "Response:\n"
	}

	endOfLine := ""
	if !strings.HasSuffix(line, "\n") {
		endOfLine = "\n"
	}
	return fmt.Sprintf("%s%s%s", responseLabel, line, endOfLine)
}

func getJsonPrintIndent() string {
	switch GetGlobalStringWithFallback(".config.restshell.json.indent", "tabs") {
	case "spaces":
		return "    "
	case "tabs":
		fallthrough
	default:
		return "\t"
	}
}
