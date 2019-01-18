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

	"golang.org/x/crypto/ssh/terminal"
)

var LastError int = 0
var InitDirectory string = ""
var ExecutableDirectory string = ""
var initialized = false

// Default parameter line for commands
var defaultParameters = "{CMD} [sub-command] [Command Options] [parameter]..."

func ReadLine() {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

// CommandProcessor -- Initiate the command interpretter using the given reader and options
func CommandProcessor(defaultPrompt string, reader io.Reader, singleStep bool, allowAbort bool) (int, bool) {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		defaultPrompt = ""
		singleStep = false
	}

	var quit bool
	var count int
	var shell = ""
	var prompt = defaultPrompt

	scanner := bufio.NewScanner(reader)
	if prompt != "" {
		fmt.Printf(prompt)
	}
	for !quit && scanner.Scan() {
		line, err := NewCommandLine(scanner.Text(), shell)
		if err != nil {
			line, _ = NewCommandLine("ERROR", "")
		}

		switch line.Command {
		case "":
		case "ERROR":
			LastError = 1
			if err == nil {
				err = errors.New("Invalid Command")
			}
			fmt.Fprintf(ErrorWriter(), "%s: %s\n", "Line Parse Error", err.Error())
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
			if len(line.ArgString) > 0 {
				shell = line.ArgString
				if !(shell[len(shell)-1] == '\\' || shell[len(shell)-1] == '/') {
					shell = shell + " "
				}
				prompt = defaultPrompt + shell
			} else {
				shell = ""
				prompt = defaultPrompt
			}
		default:
			if !line.IsComment {
				cmd, err, contStepping := processCommand(line, singleStep)
				singleStep = contStepping
				if IsFlowControl(err, FlowQuit) {
					// Flow quit is considered success; but last error remains
					quit = true
				} else if err != nil {
					LastError = 1
					fmt.Fprintf(ErrorWriter(), "%s: %s\n", line.Command, err.Error())
					if IsFlowControl(err, FlowAbort) && allowAbort {
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

func processCommand(line *Line, singleStep bool) (Command, error, bool) {

	if line.Echo || singleStep {
		fmt.Println(line.CmdLine)
	}

	cmd, tokens, err := getCmdAndArgs(line)
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

	err = processCmd(cmd, tokens, line.Echo || singleStep)
	return cmd, err, singleStep
}

func getCmdAndArgs(line *Line) (cmd Command, tokens []string, err error) {
	cmd = nil
	err = nil

	if len(line.Command) == 0 {
		return nil, tokens, errors.New("No command parsed")
	}

	// Lookup command
	if c, ok := cmdMap[line.Command]; ok {
		cmd = c
		if _, ok := cmd.(LineProcessor); ok {
			tokens = []string{line.Command, line.CmdLine}
		} else {
			tokens = line.GetCmdAndArguments()
		}
	} else {
		err = errors.New("Invalid Command '" + line.Command + "'. Try 'help'")
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
				return errors.New("Invalid sub-command: " + subCommand)
			}
			parseTokens = makeSubTokenArray(command, tokens[2:])
		}
	}

	// Setup the call to parse command options
	set := NewCmdSet()
	set.SetParameters(defaultParameters)
	InitializeCommonCmdOptions(set, CmdHelp)
	cmd.AddOptions(set)
	set.Reset()
	err := CmdParse(set, parseTokens)
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
	line, err := NewCommandLine(input, "")
	if err != nil {
		return err
	}

	if line.IsComment {
		return nil
	}

	// Special commands
	switch line.Command {
	case "REM":
		return nil
	case "HELP":
		return nil
	}

	cmd, tokens, err := getCmdAndArgs(line)
	if err != nil {
		return err
	}

	// Strip out sub command before parsing; add it back with arguments
	subCommands, hasSub := cmdSubCommands[line.Command]
	if hasSub {
		subCommand := ""
		if len(line.ArgString) > 1 && !strings.HasPrefix(line.ArgString, "-") {
			subCommand = strings.ToUpper(tokens[1])
			if !ContainsCommand(subCommand, subCommands) {
				return errors.New("Invalid sub-command: " + line.Command + " " + subCommand)
			}
			tokens = makeSubTokenArray(line.Command, tokens[2:])
		}
	}

	// Setup the call to parse command options
	set := NewCmdSet()
	InitializeCommonCmdOptions(set, CmdHelp)
	cmd.AddOptions(set)
	set.Reset()
	err = CmdParse(set, tokens)
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
