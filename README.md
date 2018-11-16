# RestShell

NOTE: This repository is still getting structured and some key features are still being developed. Future changes may break current usage patterns until settled.

## Overview
RestShell is a command line driven program to execute commands and tests against REST API's or other services. Restshell can extend its command library with custom commands to simplify interacting with specific applications or environments.

RestShell includes scripting and assertion facilities to implement repeatable and automated tests similar to unit testing frameworks. Tests can use the provided commands within scripts or developers can create specialized commands for more complex needs. Restshell commands can help with benchmarking and load testing as well.

RestShell was designed to be used in test and developer environments to help ensure developers can reproduce issues discovered in testing as well as isolating the specific functions that have issues. Best practices can be used to ensure scripts work in various environments.

> Note: Developers are encouraged to develop custom commands targeting their application needs. A custom command that executes a logon to a custom app would be easier to use then a generic post command and extracting authentication tokens. Many of the ease-of-use features are part of private commands not included in the public version. 

See the restshell-example repository to an example on building a custom Restshell.

## Build
Assuming the go environment is setup correctly, get the code:

    go get github.com/brada954/restshell

To build, open a command window at the root of the repository and execute:

    go build

To run, execute:

    restshell

## Running RestShell
The RestShell program runs as a command shell. Execute "restshell" and you will get a prompt to enter commands.

After invoking RestShell, a command can be entered with any required or optional parameters. Try these:

```
>> get --url http://api.ipify.org/?format=json  
Response:
{"ip":"123.123.123.123"}
>> assert ISSTR ip  
>> assert NOSTR ip  
ASSERT: ip: Value was an unexpected string
>> assert EQ ip xxx  
ASSERT: ip: Values not equal: xxx!=123.123.123.123
>> q  
```

Use the help command to get information on commands available:

```
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
4. Generic REST API capabilities (GET, POST)

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
Has mentioned above variables are a powerful tool for configuring parameters of commands or tests. Private commands may have some special variables it uses to perform tasks. Having variables for "secret" data is a best practice and is recommended to keep top secret data in .user files to avoid submitting to source control. Use your own descretion on test environments, etc.

Variables used by commands should have a name representing the command or a variable that may represent a  global entity to a set of commands.

By convention, variables starting with "_" should be considered reference variables; not used by commands directly. Scripts can use reference variables to initialize variables used by commands or scripts. Variables starting with "$" are considered temporary and can be cleared in bulk (see "set --clear-tmp").

### BASE Command
Making REST calls can be made easier by using the BASE command:

    >> BASE http://mysite.com/api/v1:8080
    >> GET /books
    >> GET /magazines

## Extending RestShell
RestShell is easily extendable by developing new commands to perform typical project operations. Additional commands can perform more specific actions against specific REST Api's (including saving state) as well as perform other tasks to assist managing an environment such as OS shell commands, test setup commands, adminstrative tasks, database tasks and additional assertions as helpful. The goal of extensibility is to enable developers to augment the basic tool to simplify their day to day operations or automated testing needs. Specialized commands that save, manage and re-use state can simplify a developer or testers use of the tool.

Typical commands to add may be:
1. login
2. getdata
3. postdata

This typical commands make it easy to manage the auth context between calls. The shell provides abstractions for different types of authentication to make writing the commands easier.

### Adding a command
To add a command, create a new package to hold your commands. The command package should use an init() function to register the commands with the shell (see restshell-example\commands\example). To link the new command package with RestShell, create an init2.go file (there is an example) to have an import reference to the new command package. This go mechanism will import of the command package and its init functino can register its commands with the RestShell.

Alternatively, you can fork or copy the restshell-example repository and use your own main function and build process to use RestShell as just another package. Your main function will just call shell.RunShell() to host the RestShell shell and you can initialize all the private functions you want included.

### Shell Capabilities
More to come here but:  
There are capabilities to handle authentication contexts and variables to be shared between commands. There are helpers to handle common handling of benchmarking, output generation, and command option handling.

Review the restshell command source code to learn more.

## Dependencies
This program uses:  
    http://github.com/pborman/getopt/v2 (getopt library)  
    github.com/mitchellh/mapstructure/ (structure/map mapping)  
    golang.org/x/crypto/ssh/terminal/ (terminal detection)  

## Contributors
A special thanks for the initial contributer, PeakActivity, LLC, which has allowed the program to be provided publically under the MIT License.

Other contributers:  
Brad Andersion (brada954) (maintainer)
