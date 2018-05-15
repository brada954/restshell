# RestShell

NOTE: This repository is still getting structured and some key features are still being developed. Future changes may break current usage patterns until settled.

RestShell is a command line driven program to execute commands and tests against REST API's or other systems or services. Commands are intended to simplify the developer experience exercising the system as well as create repeatable tests that can easily be used by testers and developers.

RestShell is designed to be easily extendable by developing new commands to perform typical project operations. The tool can also have commands to perform other tasks that may include admin tasks, database tasks and more if necessary. Commands may be implemented for convenience of devlopers or testers.

RestShell has a scripting and analysis facilities to implement repeatable and disectable tests. Tests can be added through specialized commands to handle complex systems, or tests can leverage the basic capabilities and more extensive scripting.

The benefit of RestShell is it is easy to build and maintain a library of tests that can be exectued in automation systems as well as easy to use in developer environments when isolating issues under debug. Best practices can be used to ensure scripts work in various environments.

Note: Developers are encouraged to develop their own commands to target the API's they are testing. Many of the ease-of-use features are part of private commands not included in this public version. Stay tuned for more work demonstrating these capabilities as the generic commands get updated.

## Build
Assuming the go environment is setup correctly, get the code:

    go get github.com/brada954/restshell

To build, open a command window at the root of the repository and execute:

    cd RestShell
    go build

To run, execute:

    RestShell

## Running RestShell
The RestShell program runs as a command shell. Execute "RestShell" and you will get a prompt to enter commands.

After invoking RestShell, a command can be entered with any required or optional parameters. Try these:

    >> get --url http://api.ipify.org/?format=json  
    >> assert ISSTR ip  
    >> assert NOSTR ip  

Use the help command to get information on commands available:
    C:\RestShell
    >> help
    Commands Available:

    Http commands:
      GET
      POST
      BMGET

    Utility related commands:
      ASSERT
      SET
      RUN
      LOG
      ENV
    >>

Most commands support a -h|--help for getting help. For example:

    >> resv --help
    Usage: resv [-dhsv] [-u value] [resvId]
    -d, --debug      Enabled debug output
    -h, --help       Display command help
    -s, --silent     Run in silent mode
    -u, --url=value  Base url for operation
    -v, --verbose    Enabled verbose output
    >>

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

.rsconfig
    A standard config file in the repository containing default configuration and reference configuration.
.rsconfig.user
    An optional file for user personal startup configuration (not committed to repository)

The repository also contains a scripts directory containing scripts for common functions and environment setup.

The root of the repository has a cluster-tests directory which is intended to contain tests intended to test the cluster.

### Best Practices
The scripting capabilities are used to easily adjust tests for different environments or test data.

For example, the .rsconfig file initializes some basic configuration but also contains some reference configuration that can be used to easily change configuration for different environments. There are scripts that can use the reference configuration to tailor the tool to work with a local deployment or a dev, staging or production environment.

Variables used by commands should have a name representing the command or a variable that may be global to all commands. Variables starting with "_" should be considered reference variables; not used by commands, but scripts can use has reference information for setting up environment data. By convention, variables starting with "$" are considered temporary and can be cleared in bulk.

Making REST calls can be made easier by using the BASE command:

    >> BASE http://mysite.com/api/v1:8080
    >> GET /books
    >> GET /magazines

## Extending commands
Developers can create new commands that may perform specialized calls, save state or re-use state between REST API calls, use variables to store outputs and inputs to make it easier to string a series of commands together without a user having to extract the data and include it in subsequent commands.

Adding a command is as simple as creating a new package to hold your commands. The command package should use an init() function to register the commands with the shell. To link the new command package with RestShell, the init.go function could have an import reference to the new commands or any new file like init2.go can be added to the root of RestShell which can cause the import of the commands. Init.go in the root of RestShell can serve as an example. Once RestShell is recompiled the commands should be available.

An example command is provided to develop new commands, but additional examples are needed to demostrate sharing state between commands as well as a feature list for better sharing techniques between the POST and GET commands.

## Assertions (Assert)
All REST commands store responses in a history buffer such that assertions can be run against the history buffer. Assertions are designed to use a simplistic XPATH-like mechanism to identify and extract a property value in a JSON response to perform validations against.

There is support to optionally test error values and Authorization JWT tokens. Extracted values can have modifiers applied to validate variations or attributes of a property value. For example, a string can be converted to its length to compare the string length to a value. See the help command for available options.

Assertions can easily be added to perform more complex validations.

## Dependencies
This program uses:
    http://github.com/pborman/getopt/v2 (getopt library)
    github.com/mitchellh/mapstructure/ (structure/map mapping)
    golang.org/x/crypto/ssh/terminal/ (terminal detection)

## Contributors
A special thanks for the initial contributer, PeakActivity, LLC, which has allowed the program to be provided publically under the MIT License.

Other contributers:
Brad Andersion (brada954) (maintainer)