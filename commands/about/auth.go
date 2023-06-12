package about

import (
	"fmt"
	"io"
)

type AuthTopic struct {
	Key         string
	Title       string
	Description string
	About       string
}

var localAuthTopic = &AuthTopic{
	Key:         "AUTH",
	Title:       "Authentication",
	Description: "Authentication Objects for decorating REST calls",
	About: `Authenticating REST API's is a redundant task for each HTTP request,
so authentication contexts are used to store authentication information for
reuse in HTTP requests.

Every REST operation has options to facilitate passing authentication details 
with a request via cookies or headers.

The LOGIN command can create a cached authentication context used with each
REST operation without having to specify authentication options.

LOGIN supports BASIC, COOKIE and HEADER mechanisms and additional simplifications for
BEARER tokens in the Authorization header.

The LOGIN command depends on the user authenticating with a service first
to obtain the token required. The user can then supply the token to the LOGIN command
for reuse with each REST operation.

    LOGIN BEARER my_auth_token

Each REST operation uses the established authentication context until it is cleared
using Login --clear or excluded by a --no-auth option.

Authentication Options with REST Operations

The following options are common to REST operations:

--basic-auth=[user][,pwd]

The basic authentication option provides user and password details with the request.
Empty user name or password values would get a prompt for the information
(passwords are hidden when typed).

--query-param=[[name,value][&name2,value2]...]

The Query Parameter Authentication adds key,value pairs to the request.

For example:
	http://xyz.com?name=value&name2=value2

The query parameter auth context will automatically escape the values in
the query string.
`,
}

func NewAuthTopic() *AuthTopic {
	return localAuthTopic
}

func (a *AuthTopic) GetKey() string {
	return a.Key
}

func (a *AuthTopic) GetTitle() string {
	return a.Title
}

func (a *AuthTopic) GetDescription() string {
	return a.Description
}

func (a *AuthTopic) WriteAbout(o io.Writer) error {
	fmt.Fprintf(o, a.About)
	return nil
}
