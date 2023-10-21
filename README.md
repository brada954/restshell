# RestShell

## Overview

RestShell is a command line shell to execute REST or other API operations against services. RestShell is for those preferring command line interfaces over clicking through Web based programs like Postman. The shell features for maintaining state, extracting data from responses, and authentication contexts simplify the user experience over curl-like commands without shared state.

With support for scripting and golang command plugins, RestShell can provide business domain operation commands as well as perform functional, benchmark and load testing with included assertion facilities.

RestShell can provide unlimited use cases for building value quickly with API services before more complex applications are fully implemented.

### Extensibliity

Developers can extend the command library with custom commands to simplify interacting with specific applications or environments. With custom commands complex REST APIs can be wrapped with simple one word commands and options and output responses in customized formats.

Custom commands can also be shared as Golang packages for reuse by the RestShell community.

To create a custom restshell with third-party packages or your own commands consult <https://github.com/brada954/restshell-example>.

### Use cases

1. Build a command to import a CSV data file into a service via an API
2. Build a command can request information via APIs and display a human readable summary or file export
    1. Health of various services in a cluster or deploymnet
    2. Collation of data representing a business entity like an Invoice, Reservation, or upcoming calandar events
    3. Query data and merge the data with a document template
3. While interacting with a domain service, the developer can easily call commands to modify data or configuration of any related service
4. Orchestrate interactions across multiple independent services

## Build

Using Golang, get the code:

```bash
go get github.com/brada954/restshell
```

To build, open a command window at the root of the repository and execute:

```bash
go build
```

To run, execute:

```bash
restshell
```

Dependent on Golang Vesion 1.14 or later

## Running RestShell

The RestShell program runs as a command shell. Execute "restshell" and you will get a prompt to enter commands.

After invoking RestShell, a command can be entered with any required or optional parameters. Try these:

```bash
>> get --url http://api.ipify.org/?format=json  
{"ip":"123.123.123.123"}
>> assert ISSTR ip  
>> assert NOSTR ip  
ASSERT: ip: Value was an unexpected string
>> assert EQ ip xxx  
ASSERT: ip: Values not equal: xxx!=123.123.123.123
>> q  
```

Use the help command to get information on commands available:

```bash
c:\go\src\github.com\brada954\restshell>restshell
>> help
restshell [COMMAND [OPTIONS]...]

restshell is a command line driven shell to execute commands and tests against
REST APIs. restshell can be invoked with arguments respresenting a
single command to execute or without arguments to run in shell mode.

General utility commands like get and post may be used against any REST API
while specialized commands may provide options for interacting with
specific APIs.

Assertion commands and a script execution engine enable this tool to run
complicated test scripts.

For more information consult the repository README.md file

Http commands:
  BASE        GET         POST        LOGIN

Benchmark commands:
  BMGET       BMPOST

Result Processing commands:
  ASSERT

Utility commands:
  REM         RUN         QUIT        LOAD        DUMP
  SET         ALIAS       DEBUG       SILENT      DIFF
  VERBOSE     DIR         CD          LOG         ENV
  SLEEP       PAUSE

Help commands:
  ABOUT       VERSION     HELP

The following are special command modifiers:

#  A comment character which needs to be the first character on the line
@  An echo character which can echo the executing command including
   expanded variables and aliases
   
>> get --help
Usage: GET [-dhsv] [--basic-auth value] [--certs] [--delete] [--head] [--headers value] [--nocert] [--noredirect] [--out-body] [--out-cookie] [--out-full] [--out-header] [--out-short] [--query-auth value] [-u value] [service route]
     --basic-auth=value
                   Use basic auth: [user][,pwd]
     --certs       Include local certs (windows may not work with system certs)
 -d, --debug       Enabled debug output
     --delete      Use HTTP DELETE method
     --head        Use HTTP HEAD method
     --headers=value
                   Set the headers [k=v]
 -h, --help        Display command help
     --nocert      Do not validate certs
     --noredirect  Do not follow redirects on requests
     --out-body    Output the response body
     --out-cookie  Output response cookies
     --out-full    Output all response data
     --out-header  Output response headers
     --out-short   Output the short response (overrides verbose)
     --query-auth=value
                   Use query param authe: [[name,value]...]
 -s, --silent      Run in silent mode
 -u, --url=value   Base url for operation
 -v, --verbose     Enabled verbose output
>>
```

RestShell has a few command line options like -d and -v which can turn on global debug output and verbose output from all commmands instead of setting the parameter on each command.

An RestShell command can be provided as arguments to RestShell when running from a shell and RestShell will execute that command and exit:
    C:\RestShell sleep 1000

## General Operations

### Key features

The RestShell program has several key features:

1. Command line shell with common options and paradigms
2. Script support (RUN)
3. Global variables including command line substitution (SET)
4. Logging (LOG)
5. Assertions or test validation mechanisms (ASSERT)
6. Generic REST API capabilities (GET, POST)

(Generic HTTP REST capabilities and assertions are continually being updated to address needs)

### Startup

When RestShell starts, it looks for two configuration files to automatically load some configuration.

- .rsconfig

   A standard config file in the repository containing default configuration and reference configuration.

- .rsconfig.user

   An optional file for user personal startup configuration (not committed to repository)

The repository also contains a test directory for testing and demonstrating scripts.

### Testing: Assertions (Assert)

All REST commands store responses in a history buffer such that assertions can be run against the history buffer. Assertions are designed to use a simplistic XPATH-like mechanism to identify and extract a property value in a JSON response to perform validations against.

There is support to optionally test error values and Authorization JWT tokens. Extracted values can have modifiers applied to validate variations or attributes of a property value. For example, a string can be converted to its length to compare the string length to a value. See the help command for available options.

Assertions can easily be added to perform more complex validations.

## Best Practices

### Scripting

The scripting capabilities are used to easily adjust tests for different environments or test data.

For example, a .rsconfig.user file can initialize basic variables and aliases for a developers environment and additional scripts can be written to alter configuration for a different environment. For example, create a usestaging.rshell script that when runs can set variables and aliases for a staging environment. When the developer runs "run usestaging" the configuration will change. A second script called "uselocal.rshell" can change the environment back to local configuration and this script can be called from the .rsconfig.user script. There are unlimited possibilities with using scripts to configure environments or test data for assertions.

### Variables

Variables are a powerful tool for configuring parameters of commands or tests. Private commands may have some special variables it uses to perform tasks. Having variables for "secret" data is a best practice and is recommended to keep top secret data in .user files to avoid submitting to source control. Use your own descretion on test environments, etc.

Variables used by commands should have a name representing the command or a variable that may represent a  global entity to a set of commands.

By convention, variables starting with "_" should be considered reference variables like "const variables". These variables should not be modified by commands. Scripts can use reference variables to initialize variables used by commands or scripts. Note: they do not have a read only implementation at this time but may in the future.

By convention, variables starting with "$" are considered temporary and can be cleared in bulk (see "set --clear-tmp").

### BASE Command

Making REST calls can be made easier by using the BASE command:

```bash
>> BASE http://mysite.com/api/v1:8080
>> GET /books
>> GET /magazines
```

Run the **ABOUT auth** command to learn about authentication contexts to use with the REST commands.

## Contributors

A special thanks for the initial contributer, PeakActivity, LLC, which has allowed the program to be provided publically under the MIT License.

Other contributers:  
Brad Andersion (brada954) (maintainer)
