package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/pborman/getopt/v2"
	"golang.org/x/crypto/ssh/terminal"
)

var LastError int = 0
var InitDirectory string = ""
var ExecutableDirectory string = ""
var initialized = false

func ReadLine() {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

func CommandProcessor(defaultPrompt string, reader io.Reader, singleStep bool, stopOnInterrupt bool) (int, bool) {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		defaultPrompt = ""
		singleStep = false
	}

	var quit bool = false
	var count int = 0
	var shell = ""
	var prompt = defaultPrompt

	getopt.SetParameters("{CMD} [sub-command] [Command Options] [parameter]...")
	scanner := bufio.NewScanner(reader)
	if prompt != "" {
		fmt.Printf(prompt)
	}
	for !quit && scanner.Scan() {
		echo := false
		input := scanner.Text()
		input = strings.TrimSpace(input)

		// Get first token for special handling
		args := strings.SplitN(input, " ", 2)
		command := ""
		if len(args) > 0 {
			command = strings.ToUpper(args[0])
		}

		if strings.HasPrefix(command, "#") {
			command = ""
		}

		if strings.HasPrefix(command, "@") {
			echo = true
			command = strings.TrimLeft(command, "@")
			input = strings.TrimLeft(input, "@")
		}

		if len(shell) > 0 {
			input = shell + input
		}

		// Allow aliases
		if alias, err := GetAlias(command); err == nil {
			if IsDebugEnabled() {
				fmt.Fprintf(ConsoleWriter(), "Using alias: %s\n", alias)
			}
			input = alias
			if len(args) > 1 {
				input = input + " " + args[1]
			}
			tmp := strings.SplitN(input, " ", 2)
			command = strings.ToUpper(tmp[0])
		}

		switch command {
		case "":
		case "QUIT":
			fallthrough
		case "Q":
			if len(shell) == 0 {
				quit = true
			} else {
				prompt = defaultPrompt
				shell = ""
			}
		case "SHELL":
			if !quit && len(args) == 2 && len(args[1]) > 0 {
				shell = args[1]
				if !(shell[len(shell)-1] == '\\' || shell[len(shell)-1] == '/') {
					shell = shell + " "
				}
				prompt = defaultPrompt + shell
				input = ""
			} else {
				shell = ""
				prompt = defaultPrompt
			}
		default:
			if !strings.HasPrefix(input, "#") {
				cmd, err, contStepping := processCommand(input, echo, singleStep)
				singleStep = contStepping
				if IsFlowControl(err, FlowQuit) {
					quit = true
				} else if err != nil {
					LastError = 1
					fmt.Fprintf(ErrorWriter(), "%s: %s\n", command, err.Error())
					if IsFlowControl(err, FlowAbort) && stopOnInterrupt {
						quit = true
					}
				} else if track, trackable := cmd.(Trackable); cmd != nil && trackable {
					if track.DoNotCount() == false {
						count++
					}
					if track.DoNotClearError() == false {
						LastError = 0
					}
					count = count + track.CommandCount()
				} else {
					LastError = 0
					count++
				}

				if CommandRequestsQuit(cmd) {
					quit = true
				}
			}
		}
		if !quit && prompt != "" {
			fmt.Printf(prompt)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(ErrorWriter(), "Scanner error %s\n", err.Error())
		return count, false
	}
	return count, true
}

func GetInitDirectory() string {
	return InitDirectory
}

func GetExeDirectory() string {
	return ExecutableDirectory
}

func InitializeShell() {
	if initialized == true {
		return
	}

	initialized = true
	curdir, err := os.Getwd()
	if err == nil {
		InitDirectory = curdir
		if len(curdir) > 0 && strings.HasSuffix(curdir, "/") == false {
			InitDirectory = InitDirectory + "/"
		}
	}

	exPath := ""
	{
		ex, err := os.Executable()
		if err == nil {
			exPath = filepath.Dir(ex)
		}
	}
	ExecutableDirectory = exPath
}

func processCommand(line string, echo bool, singleStep bool) (Command, error, bool) {
	echoed := echo || singleStep
	cmd, tokens, err := getCmdAndArgs(line, echoed)
	if err != nil {
		return cmd, err, singleStep
	}

	if singleStep {
		switch doSingleStep(cmd) {
		case "q":
			return cmd, NewFlowError("Quit requested", FlowQuit), singleStep
		case "g":
			singleStep = false
		default:
		}
	}

	err = processCmd(cmd, tokens, echoed)
	return cmd, err, singleStep
}

func getCmdAndArgs(input string, echo bool) (cmd Command, tokens []string, err error) {
	cmd = nil
	err = nil

	input = PerformVariableSubstitution(input)
	tokens = LineParse(input)

	if len(tokens) < 1 {
		return nil, tokens, errors.New("Parse failed to find tokens")
	}

	var command = strings.ToUpper(tokens[0])

	// Lookup command
	if c, ok := cmdMap[command]; ok {
		cmd = c
		if _, ok := cmd.(LineProcessor); ok {
			tokens = []string{command, input}
		} else {
			tokens[0] = command
		}
	} else {
		err = errors.New("Invalid Command '" + command + "'. Try 'help'")
	}

	if echo {
		fmt.Println(input)
	}
	return
}

func doSingleStep(cmd Command) string {
	// Single step skips over commands like REM and other non-counting commands
	if CommandRequestsNoStep(cmd) == false {
		return getStepCommand()
	}
	return ""
}

func parseAndExecute(cmd Command, command string, tokens []string) error {
	// Strip out sub command before parsing; add it back with arguments
	parseTokens := tokens
	subCommands, hasSub := cmdSubCommands[command]
	subCommand := ""
	if hasSub {
		if len(tokens) > 1 && !strings.HasPrefix(tokens[1], "-") {
			subCommand = strings.ToUpper(tokens[1])
			if !ContainsCommand(subCommand, subCommands) {
				return errors.New("Invalid sub-command: " + command + " " + subCommand)
			}
			parseTokens = makeSubTokenArray(command, tokens[2:])
		}
	}

	// Setup the call to parse command options
	set := getopt.New()
	InitializeCommonCmdOptions(set, CmdHelp)
	cmd.AddOptions(set)
	set.Reset()
	err := set.Getopt(parseTokens, nil)
	if err != nil {
		fmt.Fprintln(ErrorWriter(), err.Error())
		set.Usage()
		return errors.New("Invalid arguments")
	}
	if IsCmdHelpEnabled() {
		set.Usage()
		return nil
	}

	if hasSub && subCommand != "" {
		return cmd.Execute(makeSubTokenArray(subCommand, set.Args()))
	} else {
		return cmd.Execute(set.Args())
	}
}

func processCmd(cmd Command, tokens []string, echoed bool) (result error) {
	command := tokens[0]
	result = nil

	if command == "HELP" {
		DisplayHelp()
		return nil
	}

	if cmd == nil || len(tokens) < 1 {
		return errors.New("Failed to process command line")
	}

	// Setup interrupt handler
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	var interrupted bool = false
	defer func() {
		signal.Stop(sigchan)
		if interrupted {
			result = NewFlowError("Command interrupted", FlowAbort)
			_ = recover()
		} else if r := recover(); r != nil {
			result = errors.New("Command failed")
			message := fmt.Sprintf("Exception processing %s command", command)
			fmt.Fprintln(ErrorWriter(), message)
			fmt.Fprintf(ErrorWriter(), "Panic: %v\n%s\n", r, debug.Stack())
		}
	}()

	// Process command after setting up interrupt handler
	go func() {
		for sig := range sigchan {
			if sig == os.Interrupt {
				interrupted = true
				if abort, ok := cmd.(Abortable); ok {
					abort.Abort()
				}
			}
		}
	}()

	if CommandProcessesLine(cmd) == true {
		// Line processors do not have standard argument processing
		if c, ok := cmd.(LineProcessor); ok {
			result = c.ExecuteLine(tokens[1], echoed)
		}
		return result
	}
	return parseAndExecute(cmd, command, tokens)
}

func validateCmd(input string) error {
	input = PerformVariableSubstitution(input)
	var tokens = LineParse(input)
	var command = strings.ToUpper(tokens[0])

	// Special commands
	switch command {
	case "REM":
		return nil
	case "HELP":
		return nil
	}

	var cmdServer Command
	{
		var ok bool
		cmdServer, ok = cmdMap[command]
		if !ok {
			return errors.New("Invalid Command")
		}
	}

	// Strip out sub command before parsing; add it back with arguments
	subCommands, hasSub := cmdSubCommands[command]
	subCommand := ""
	if hasSub {
		if len(tokens) > 1 && !strings.HasPrefix(tokens[1], "-") {
			subCommand = strings.ToUpper(tokens[1])
			if !ContainsCommand(subCommand, subCommands) {
				return errors.New("Invalid sub-command: " + command + " " + subCommand)
			}
			tokens = makeSubTokenArray(command, tokens[2:])
		}
	}

	// Setup the call to parse command options
	set := getopt.New()
	InitializeCommonCmdOptions(set, CmdHelp)
	cmdServer.AddOptions(set)
	set.Reset()
	err := set.Getopt(tokens, nil)
	if err != nil {
		return errors.New("Invalid arguments")
	}
	return nil
}

func makeSubTokenArray(subCmd string, tokens []string) []string {
	result := []string{}
	result = append(result, subCmd)
	result = append(result, tokens...)
	return result
}

func ContainsCommand(cmd string, tokens []string) bool {
	for _, v := range tokens {
		if v == cmd {
			return true
		}
	}
	return false
}

func getPassword(prompt string) string {
	if len(prompt) > 0 {
		fmt.Fprintf(os.Stdout, prompt)
	}

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Fprintf(os.Stdout, "\n")
	if err != nil {
		return ""
	}
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

func getLine(prompt string) string {
	if len(prompt) > 0 {
		fmt.Fprintf(os.Stdout, prompt)
	}

	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		input := scanner.Text()
		input = strings.Trim(input, " ")
		return input
	}
	return ""
}

// Return values for continue stepping and quit or not
func getStepCommand() string {
	var stepCmd string
	fmt.Print("Stopped> ")
	fmt.Scanln(&stepCmd)
	stepCmd = strings.ToLower(strings.TrimSpace(stepCmd))
	if len(stepCmd) == 0 {
		return ""
	}
	if stepCmd == "g" {
		return "g"
	} else if stepCmd == "q" {
		return "q"
	}
	return ""
}
