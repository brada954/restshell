package rest

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/brada954/restshell/shell"
)

type HeaderAuth struct {
	headers http.Header
}

// NewHeaderAuth -- Create an Auth context that adds cookies to the request
func NewHeaderAuth() *HeaderAuth {
	headers := make(http.Header, 0)
	return &HeaderAuth{headers: headers}
}

func (a *HeaderAuth) IsAuthed() bool {
	return len(a.headers) > 0
}

func (a *HeaderAuth) AddAuth(req *http.Request) {
	if a.IsAuthed() {
		for k, data := range a.headers {
			for _, v := range data {
				if shell.IsCmdDebugEnabled() {
					fmt.Fprintf(shell.ConsoleWriter(), "Adding header to request: %s=%s\n", k, v)
				}
				req.Header.Add(k, v)
			}
		}
	} else {
		if shell.IsCmdDebugEnabled() {
			fmt.Fprintln(shell.ConsoleWriter(), "HeaderAuth is not authed, no cookies to add")
		}
	}
}

func (a *HeaderAuth) ToString() string {
	var buf bytes.Buffer
	var separator = ""
	for k, data := range a.headers {
		fmt.Fprintf(&buf, "%s%s=", separator, k)
		var comma = ""
		for _, v := range data {
			fmt.Fprintf(&buf, "%s%s", comma, v)
			comma = ","
		}
		separator = "; "
	}
	return buf.String()
}

func (a *HeaderAuth) AddHeader(name string, value string) {
	if data, ok := a.headers[name]; ok {
		a.headers[name] = append(data, value)
	} else {
		data := make([]string, 0)
		data = append(data, value)
		a.headers[name] = data
	}
}
