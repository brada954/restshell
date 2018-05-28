package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type RunCommand struct {
	running         int
	interrupted     bool
	header          *bool
	list            *bool
	ifCondition     *string
	stepOption      *bool
	iterationOption *int
	execOption      *bool
	// Note: nesting makes count only valid from end of execute and calling CommandCount() immediately
	count int
}

var DefaultScriptExtension = ".rshell"

func NewRunCommand() *RunCommand {
	cmd := RunCommand{running: 0}
	return &cmd
}

func (cmd *RunCommand) AddOptions(set CmdSet) {
	set.SetParameters("scripts...")
	cmd.ifCondition = set.StringLong("cond", 0, "", "run script if specified variable is not empty")
	cmd.list = set.BoolLong("list", 0, "List the contexts of script file")
	cmd.header = set.BoolLong("header", 0, "Display header of script file (Leading REM commands)")
	cmd.stepOption = set.BoolLong("step", 0, "Single step through script")
	cmd.iterationOption = set.IntLong("iterations", 'i', 1, "run the script iteration number of times")
	cmd.execOption = set.BoolLong("exec", 0, "Execute quoted parameters as script commands")
	AddCommonCmdOptions(set, CmdDebug, CmdVerbose, CmdSilent)
}

// Validate the script exists by either basename or basename plus suffix
// return the file name modified with extension if necesssary.
// If the error is not a file existence problem the file is returned.
func ValidateScriptExists(file string) (string, error) {
	if len(file) == 0 {
		return "", errors.New("The file does not exist")
	}
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		if !strings.HasSuffix(strings.ToLower(file), DefaultScriptExtension) {
			if _, err2 := os.Stat(file + DefaultScriptExtension); err2 == nil {
				file = file + DefaultScriptExtension
				err = err2
				if IsCmdDebugEnabled() {
					fmt.Fprintf(ConsoleWriter(), "Appending extension to file name: %s\n", file)
				}
			}
		}

		if os.IsNotExist(err) { // Still not exists
			if IsCmdDebugEnabled() {
				fmt.Fprintf(ConsoleWriter(), "Unable to open file: %s\n", file)
			}
			return "", errors.New("The file does not exist")
		}
	}

	if err != nil {
		return file, errors.New("Error accessing file")
	}
	return file, nil
}

func (cmd *RunCommand) executeFile(file string, runSilent bool) (count int, elapsed time.Duration, result error) {
	count = 0
	elapsed = 0

	if verifyCondition(*cmd.ifCondition) == false {
		fmt.Fprintf(OutputWriter(), "Missing Condition: %s. Skipping script: %s\n", *cmd.ifCondition, file)
		return 0, 0, nil
	}

	var path = ""
	{
		scriptFile, err := ValidateScriptExists(file)
		if err != nil {
			if runSilent {
				// We do not care about file existance issues
				return 0, 0, nil
			}
			return 0, 0, err
		}

		if abspath, err := filepath.Abs(scriptFile); err == nil {
			path = filepath.Dir(abspath)
		}
		file = scriptFile
	}

	if IsCmdDebugEnabled() || IsCmdVerboseEnabled() {
		fmt.Fprintf(ConsoleWriter(), "Processing file: %s\n", file)
	}
	h, err := os.Open(file)
	if err != nil {
		return 0, 0, errors.New("Failed to read script file: " + err.Error())
	}

	curdir, err := os.Getwd()
	if err != nil {
		curdir = ""
	}

	cmd.running = cmd.running + 1
	defer func(f *os.File, c *RunCommand) {
		f.Close()
		c.running = c.running - 1
		if c.running < 0 {
			c.running = 0
		}
		if len(curdir) > 0 {
			if IsCmdDebugEnabled() {
				fmt.Fprintln(ConsoleWriter(), "RUN resetting working directory: ", curdir)
			}
			os.Chdir(curdir)
		}
	}(h, cmd)

	if len(curdir) > 0 && len(path) > 0 {
		if IsCmdDebugEnabled() {
			fmt.Fprintln(ConsoleWriter(), "RUN setting working directory: ", path)
		}
		os.Chdir(path)
	}

	if *cmd.header || *cmd.list {
		listfile(h, *cmd.header)
		return 0, 0, nil
	}

	startTime := time.Now()
	commands, success := CommandProcessor("", h, *cmd.stepOption, true)
	elapsed = time.Since(startTime)
	if !success {
		return commands, elapsed, errors.New("Command processor failed")
	}
	return commands, elapsed, nil
}

func (cmd *RunCommand) executeStream(r io.Reader, runSilent bool) (count int, elapsed time.Duration, result error) {
	count = 0
	elapsed = 0

	cmd.running = cmd.running + 1
	defer func(c *RunCommand) {
		c.running = c.running - 1
		if c.running < 0 {
			c.running = 0
		}
	}(cmd)

	if *cmd.header || *cmd.list {
		listfile(r, *cmd.header)
		return 0, 0, nil
	}

	startTime := time.Now()
	commands, success := CommandProcessor("", r, *cmd.stepOption, true)
	elapsed = time.Since(startTime)
	if !success {
		return commands, elapsed, errors.New("Command processor failed")
	}
	return commands, elapsed, nil
}

func (cmd *RunCommand) Execute(args []string) error {
	// Cache silent config as it changes with script running
	cmd.count = 0
	cmd.interrupted = false
	iterations := *cmd.iterationOption

	if len(args) == 0 {
		return errors.New("Need to specify at least one file to run")
	}

	if cmd.running > 3 {
		return errors.New("Too many nested scripts script")
	}

	runSilent := IsCmdSilentEnabled()
	i := iterations
	var result error = nil
	var commands int
	var duration time.Duration
	resultMsg := "Success"
	for ; i > 0; i-- {
		if *cmd.execOption {
			str := strings.Join(args, "\n")
			r := strings.NewReader(str)
			count, elapsed, err := cmd.executeStream(r, runSilent)
			commands = commands + count
			duration = duration + elapsed
			if err != nil {
				result = err
				resultMsg = err.Error()
				break
			}
		} else {
			for _, fileName := range args {
				count, elapsed, err := cmd.executeFile(fileName, runSilent)
				commands = commands + count
				duration = duration + elapsed
				if err != nil {
					if len(args) > 1 {
						fmt.Fprintf(ErrorWriter(), "Aborting due to errors in script: %s\n", fileName)
					}
					// TBD: HOw fatal should "not exists" be considered in list of files
					return err
				}
			}
		}
	}

	cmd.count = commands
	if result == nil && LastError != 0 {
		resultMsg = "Last command returned an error"
	}

	if !runSilent || IsCmdDebugEnabled() {
		if iterations > 1 {
			fmt.Fprintf(OutputWriter(), "Ran %d commands in %s over %d iterations. Exited with %s\n",
				commands, getDurationString(duration), iterations, resultMsg)
		} else {
			fmt.Fprintf(OutputWriter(), "Ran %d commands in %s. Exited with %s\n",
				commands, getDurationString(duration), resultMsg)
		}
	}
	return nil
}

func (cmd *RunCommand) DoNotCount() bool {
	return true
}

func (cmd *RunCommand) DoNotClearError() bool {
	return true
}

func (cmd *RunCommand) CommandCount() int {
	return cmd.count
}

func (cmd *RunCommand) Abort() {
	cmd.interrupted = true
}

func getDurationString(duration time.Duration) string {
	if duration > 2*time.Second {
		return strconv.FormatFloat(float64(duration)/float64(time.Second), 'f', 3, 64) + "s"
	} else {
		return strconv.FormatFloat(float64(duration)/float64(time.Millisecond), 'f', 1, 64) + "ms"
	}
}

func verifyCondition(variable string) bool {
	variable = strings.TrimSpace(variable)
	if len(variable) <= 0 {
		return true
	}

	parts := strings.Split(variable, "=")
	variable = parts[0]

	value := GetGlobal(variable)
	if str, ok := value.(string); ok {
		if len(strings.TrimSpace(str)) > 0 {
			if len(parts) > 1 {
				if str == parts[1] {
					return true
				}
			} else {
				return true
			}
		}
	} else if value != nil {
		return true
	}
	return false
}

func listfile(reader io.Reader, onlyHeader bool) {
	scanner := bufio.NewScanner(reader)
	quit := false
	for !quit && scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimSpace(input)

		// Get first token for special handling
		args := strings.SplitN(input, " ", 2)
		command := ""
		if len(args) > 0 {
			command = strings.ToUpper(args[0])
		}

		if strings.HasPrefix(command, "REM") || onlyHeader == false {
			fmt.Fprintln(OutputWriter(), input)
		} else {
			quit = true
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(ErrorWriter(), "Scanner error %s\n", err.Error())
	}
}
