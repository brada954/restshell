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
