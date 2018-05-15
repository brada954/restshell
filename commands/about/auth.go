package about

type AuthTopic struct {
	Key string
    Title string
    Description string
    About string
}

var localAuthTopic = &AuthTopic{
	Key: "AUTH",
	Title: "Authentication",
	Description: "Authentication Objects for decorating REST calls",
	About:
`Authenticating REST API's is a rundandant task that would be cumbersome to
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
    - NoAuth

JWT
JWT is a special authentication that requires a request to get the token.
There will be a separate command to perform a JWT authentication that will
store the token for future use by other commands. An authentication command
has a well known identifier other REST commands to use. Examples are:

    MSAUTH resvno lastname

This command generates a token that RESV and BUY commands use. Commands
may provide options to not include authentication or override authentication
with a different mechanism. When the command uses JWT, the proper 
Authorization Header is updated in the REST API request. The RESV and BUY
commands may fail if MSAUTH is not called first.

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
added has query parameters on the request. The "MOR" command is an
example that has a "LOGIN" and "KEYS" sub-command to initialize the
authenticadtion object for all other MOR sub-commands to be a Basic
Authentication Object or a Query Parameter Authenticadtion object.

    Option syntax;
    --query-param=[[name,value]...]

The name value is used as the query parameter name and the value is
used as the value of the query parameter. For example:

    http://xyz.com?name=value&name2=value2
	
The query parameter object will automatically escape the values in
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

func (a *AuthTopic) GetAbout() string {
	return a.About
}