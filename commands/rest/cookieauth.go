package rest

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/brada954/restshell/shell"
)

type CookieAuth struct {
	cookies []*http.Cookie
}

// NewCookieAuth -- Create an Auth context that adds cookies to the request
func NewCookieAuth() *CookieAuth {
	cookies := make([]*http.Cookie, 0)
	return &CookieAuth{cookies: cookies}
}

func (a *CookieAuth) IsAuthed() bool {
	return len(a.cookies) > 0
}

func (a *CookieAuth) AddAuth(req *http.Request) {
	if a.IsAuthed() {
		for _, c := range a.cookies {
			if shell.IsCmdDebugEnabled() {
				fmt.Fprintf(shell.ConsoleWriter(), "Adding cookie to request: %s=%s\n", c.Name, c.Value)
			}
			req.AddCookie(c)
		}
	} else {
		if shell.IsCmdDebugEnabled() {
			fmt.Fprintln(shell.ConsoleWriter(), "CookieAuth is not authed, no cookies to add")
		}
	}
}

func (a *CookieAuth) ToString() string {
	var buf bytes.Buffer
	var separator = ""
	for _, c := range a.cookies {
		fmt.Fprintf(&buf, "%s%s=%s", separator, c.Name, c.Value)
		separator = "; "
	}
	return buf.String()
}

func (a *CookieAuth) AddCookie(name string, value string) {
	a.cookies = append(a.cookies, &http.Cookie{Name: name, Value: value})
}
