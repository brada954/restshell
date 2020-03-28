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
	About: `Authenticating REST API's is a rundandant task that would be cumbersome to
perform with each HTTP request, so authentication objects are used to store
authentication information that can be applied to HTTP requests. There may 
be various well known objects stored in the global store that can be
configured and used by different commands. For example, gateway API's may
use the common JWT authentication token, but calls to back end services
may use other mechanisms like query parameters or Basic Authentication

Since HTTP requests support multiple methods of authentication like cookies, 
Authorization headers for basic authentiation or JWTs. There are different
types of authentication objects. These types of objects may include:

    - JWT
    - Basic Authentication
	- Query Parameters
	- Cookies
	- Headers
    - NoAuth

JWT
JWT is a special authentication that requires a request to get the token.
A custom command can be created to authenticate a user to a service and
save the token for future use by other commands.

Basic Authentication
Basic Authentication leverages a username and password pair that gets
encoded and passed on the REST api. Commands may exist to create an object
stored in the global store or a REST api may accept an option like 
--basic-auth which enables a user to provide a user name and password.

    Option syntax:
    --basic-auth=[user][,pwd]

For basic authentication an empty user name or password would get a prompt
for the information.

Query Parameter Authentication
Query Parameteer authentication enables a set of key,value pairs to be
added has query parameters on the request.

    Option syntax;
    --query-param=[[name,value]...]

The name value is used as the query parameter name and the value is
used as the value of the query parameter. For example:

    http://xyz.com?name=value&name2=value2
	
The query parameter object will automatically escape the values in
the query string.

Cookie Authentication
Cookie authentication enables cookies to be defined that suffice
as authentication. 

Header Authentication 
Header authentication allows headers to be set to specific values
for authentication.

There is a LOGIN command capable of creating auth contexts based on
cookies or headers. The basic REST commands work with the AuthContext
created by LOGIN. Additional custom commands can use those contexts
as well.
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
