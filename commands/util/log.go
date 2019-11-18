package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/brada954/restshell/shell"
)

type LogCommand struct {
	cmdTruncate *bool
	cmdAppend   *bool
	cmdStop     *bool
	logFile     *os.File
}

func NewLogCommand() *LogCommand {
	cmd := LogCommand{logFile: nil}
	return &cmd
}

func (cmd *LogCommand) AddOptions(set shell.CmdSet) {
	cmd.cmdTruncate = set.BoolLong("truncate", 0, "Truncate the log file first")
	cmd.cmdAppend = set.BoolLong("append", 'a', "Append to an existing file")
	cmd.cmdStop = set.BoolLong("stop", 0, "Stop the current log")
}

func (cmd *LogCommand) Execute(args []string) error {
	if *cmd.cmdStop {
		return performStop(cmd, args)
	} else {
		return performStart(cmd, args)
	}
}

func performStop(cmd *LogCommand, args []string) error {
	if len(args) > 0 {
		fmt.Fprintf(shell.ErrorWriter(), "Arguments not allowed with stop option")
	} else if *cmd.cmdTruncate {
		fmt.Fprintf(shell.ErrorWriter(), "Truncate option ignored with --stop option")
	}

	if cmd.logFile == nil {
		return errors.New("Logging not currently enabled")
	}

	_, err := shell.ResetOutput()

	cmd.logFile.Close()
	cmd.logFile = nil
	return err
}

func performStart(cmd *LogCommand, args []string) error {
	if len(args) == 0 {
		return errors.New("File name is required argument")
	}
	var name = args[0]
	var file *os.File
	if _, err := os.Stat(name); err == nil {
		if !(*cmd.cmdTruncate || *cmd.cmdAppend) {
			return errors.New("File exists; use --append or --truncate to use the file")
		}
		flags := os.O_APPEND | os.O_WRONLY
		if *cmd.cmdTruncate {
			flags = os.O_WRONLY
		}
		file, err = os.OpenFile(name, flags, 0644)
		if err != nil {
			return errors.New("Open failed: " + err.Error())
		}
		if *cmd.cmdTruncate {
			file.Truncate(0)
		}
	} else {
		file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New("Open failed: " + err.Error())
		}
	}

	fmt.Fprintf(shell.ConsoleWriter(), "Opened logging...\n")
	cmd.logFile = file
	return shell.SetOutput(file)
}
