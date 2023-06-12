package rest

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/brada954/restshell/shell"
)

type LoginCommand struct {
	// Place getopt option value pointers here
	optionClear *bool
}

func NewLoginCommand() *LoginCommand {
	return &LoginCommand{}
}

func (cmd *LoginCommand) GetSubCommands() []string {
	var commands = []string{"COOKIE", "HEADER", "BEARER", "BASIC"}
	return shell.SortedStringSlice(commands)
}

func (cmd *LoginCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("logintype [([name=value]... [name=value;name=value;...]|[token])]")
	cmd.optionClear = set.BoolLong("clear", 0, "Clear the auth context")

	set.SetUsage(func() {
		cmd.HeaderUsage(shell.ConsoleWriter())
		set.PrintUsage(shell.ConsoleWriter())
		cmd.ExtendedUsage(shell.ConsoleWriter())
	})

	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *LoginCommand) HeaderUsage(w io.Writer) {
	fmt.Fprintln(w, "LOGIN logintype [parameters...]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, `Create an authentication context based on the desired type`)
	fmt.Fprintln(w)
}

// ExtendedUsage -- write the extended useage
func (cmd *LoginCommand) ExtendedUsage(w io.Writer) {
	fmt.Fprintf(w, "\nAuthentication Types\n")
	lines := shell.ColumnizeTokens(cmd.GetSubCommands(), 4, 15)
	for _, v := range lines {
		fmt.Fprintf(w, "  %s\n", v)
	}
}

// Execute - execute the given command
func (cmd *LoginCommand) Execute(args []string) error {
	// Validate arguments
	if *cmd.optionClear {
		shell.SetAuthContext(RESTBASEAUTHKEY, nil)
		if len(args) > 0 {
			fmt.Fprintln(shell.ConsoleWriter(), "WARNING: Login context was cleared; parameters ignored")
		}
		return nil
	}

	if len(args) < 1 {
		return shell.ErrInvalidSubCommand
	}

	switch args[0] {
	case "COOKIE":
		return cmd.setCookieAuth(args[1:])
	case "HEADER":
		return cmd.setHeaderAuth(args[1:])
	case "BEARER":
		return cmd.setBearerAuth(args[1:])
	case "BASIC":
		return cmd.setBasicAuth(args[1:])
	default:
		return shell.ErrInvalidSubCommand
	}
}

func (cmd *LoginCommand) setCookieAuth(args []string) error {
	authContext := NewCookieAuth()
	errCnt := 0

	if len(args) == 0 {
		line := shell.GetLine("Input Cookie Data [x=y;a=b]:\n")
		args = shell.LineParse(line)
	}

	// Execute commands
	for _, arg := range args {
		cookies := strings.Split(arg, ";")
		for _, c := range cookies {
			pair := strings.SplitN(c, "=", 2)
			if len(pair) == 2 {
				authContext.AddCookie(strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1]))
			} else {
				fmt.Fprintf(shell.ErrorWriter(), "Skipping invalid cookie: %s\n", c)
				errCnt++
			}
		}
	}

	if errCnt == 0 {
		shell.SetAuthContext(RESTBASEAUTHKEY, authContext)
		return nil
	}
	return errors.New("Invalid parameters, no authentication saved")
}

func (cmd *LoginCommand) setHeaderAuth(args []string) error {
	authContext := NewHeaderAuth()
	errCnt := 0

	if len(args) == 0 {
		line := shell.GetLine("Input Header Data [x=y;a=b,c]:\n")
		args = shell.LineParse(line)
	}

	// Execute commands
	for _, arg := range args {
		headers := strings.Split(arg, ";")
		for _, c := range headers {
			pair := strings.SplitN(c, "=", 2)
			if len(pair) == 2 {
				for _, v := range strings.Split(pair[1], ",") {
					authContext.AddHeader(strings.TrimSpace(pair[0]), strings.TrimSpace(v))
				}
			} else {
				fmt.Fprintf(shell.ErrorWriter(), "Skipping invalid header: %s\n", c)
				errCnt++
			}
		}
	}

	if errCnt == 0 {
		shell.SetAuthContext(RESTBASEAUTHKEY, authContext)
		return nil
	}
	return errors.New("Invalid parameters, no authentication saved")
}

func (cmd *LoginCommand) setBearerAuth(args []string) error {
	if len(args) == 0 {
		line := shell.GetLine("Input Bearer Token value:\n")
		args = shell.LineParse(line)
	}

	if len(args) > 1 {
		return errors.New("too many arguments provided, only one token value supported")
	}

	// Execute commands
	authContext := NewHeaderAuth()
	authContext.AddHeader("Authorization", "Bearer "+strings.TrimSpace(args[0]))
	shell.SetAuthContext(RESTBASEAUTHKEY, authContext)
	return nil
}

func (cmd *LoginCommand) setBasicAuth(args []string) error {
	if len(args) > 2 {
		return errors.New("too many arguments provided, only one token value supported")
	}

	u := ""
	{
		if len(args) > 0 {
			u = args[0]
		}
	}
	p := ""
	{
		if len(args) > 1 {
			u = args[1]
		}
	}

	// Execute commands
	authContext := shell.NewBasicAuth(u, p)
	shell.SetAuthContext(RESTBASEAUTHKEY, authContext)
	return nil
}
