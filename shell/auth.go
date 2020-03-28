package shell

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type QueryParmKeyValue struct {
	Key   string
	Value string
}

type Auth interface {
	IsAuthed() bool
	AddAuth(*http.Request)
	ToString() string
}

type BasicAuth struct {
	UserName string
	Password string
}

type CookieAuth struct {
	CookieName string
	AuthToken  string
}

type JwtHeaderAuth struct {
	AuthToken string
}

type NoAuth struct {
}

type QueryParamAuth struct {
	KeyPairs []QueryParmKeyValue
}

var AnonymousAuth = NoAuth{}

var authContexts map[string]Auth = map[string]Auth{
	"Anon":    AnonymousAuth,
	"Default": BasicAuth{},
}

func GetAuthContext(ctx string) (Auth, error) {
	auth, ok := authContexts[ctx]
	if !ok || auth == nil {
		return nil, errors.New("Not Found")
	}
	return auth, nil
}

func SetAuthContext(ctx string, auth Auth) {
	authContexts[ctx] = auth
}

func (a BasicAuth) IsAuthed() bool {
	return len(a.UserName) > 0
}

func (a BasicAuth) AddAuth(req *http.Request) {
	req.SetBasicAuth(a.UserName, a.Password)
}

func (a BasicAuth) ToString() string {
	return a.UserName
}

func NewBasicAuth(u string, p string) BasicAuth {
	if len(u) == 0 {
		u = GetLine("Username: ")
	}
	if len(p) == 0 {
		p = GetPassword("Password: ")
	}
	return BasicAuth{UserName: u, Password: p}
}

func (a NoAuth) IsAuthed() bool {
	return false
}

func (a NoAuth) AddAuth(req *http.Request) {
	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "AddAuth called on anonymous user")
	}
}

func (a NoAuth) ToString() string {
	return "{no auth}"
}

func NewJwtHeaderAuth(t string) JwtHeaderAuth {
	return JwtHeaderAuth{t}
}

func (a JwtHeaderAuth) IsAuthed() bool {
	return a.AuthToken != ""
}

func (a JwtHeaderAuth) AddAuth(req *http.Request) {
	if a.IsAuthed() {
		if IsCmdDebugEnabled() {
			fmt.Fprintf(ConsoleWriter(), "Adding Auth Header to request: %s\n", a.AuthToken)
		}
		req.Header.Set("Authorization", "Bearer "+a.AuthToken)
	} else {
		if IsCmdDebugEnabled() {
			fmt.Fprintln(ConsoleWriter(), "JWT missing token, not adding")
		}
	}
}

func (a JwtHeaderAuth) ToString() string {
	return a.AuthToken[0:15] + "..."
}

func NewQueryParamAuth(kv ...string) QueryParamAuth {
	auth := QueryParamAuth{}

	auth.KeyPairs = make([]QueryParmKeyValue, 0)
	for i := 0; i < len(kv)-1; i = i + 2 {
		key := strings.TrimSpace(kv[i])
		value := strings.TrimSpace(kv[i+1])
		auth.KeyPairs = append(auth.KeyPairs, QueryParmKeyValue{key, value})
	}
	return auth
}

func (a QueryParamAuth) IsAuthed() bool {
	return a.KeyPairs != nil && len(a.KeyPairs) > 0
}

func (a QueryParamAuth) AddAuth(req *http.Request) {
	params := ""
	sep := ""
	for _, nvp := range a.KeyPairs {
		params = params + sep + nvp.Key + "=" + url.QueryEscape(nvp.Value)
		sep = "&"
	}

	if params == "" {
		return
	}

	newurl := req.URL.String()
	if strings.Contains(newurl, "?") {
		newurl = newurl + "&" + params
	} else {
		newurl = newurl + "?" + params
	}

	if result, err := url.Parse(newurl); err == nil {
		req.URL = result
	}
}

func (a QueryParamAuth) ToString() string {
	keys := make([]string, 0)
	for _, pair := range a.KeyPairs {
		keys = append(keys, pair.Key)
	}
	return strings.Join(keys, ",")
}

func (a QueryParamAuth) GetKeyValue(key string) (string, bool) {
	for _, nvp := range a.KeyPairs {
		if strings.ToLower(nvp.Key) == strings.ToLower(key) {
			return nvp.Value, true
		}
	}
	return "", false
}
