package util

import (
	"errors"
	"log"
	"os/exec"
	"time"

	"github.com/brada954/restshell/shell"
)

var ()

type DiffCommand struct {
	// Globals

	// Options
}

func NewDiffCommand() *DiffCommand {
	return &DiffCommand{}
}

func (cmd *DiffCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("[-- {git options for git diff --no-index}] file1 file2")
	set.SetUsage(func() {
		set.PrintUsage(shell.ConsoleWriter())
	})
	shell.AddCommonCmdOptions(set, shell.CmdDebug)
}

func (cmd *DiffCommand) Execute(args []string) error {

	gitoptions := []string{"diff", "--no-index"}
	gitoptions = append(gitoptions, args...)
	git := exec.Command("git", gitoptions...)

	git.Stdout = shell.OutputWriter()
	git.Stderr = shell.ErrorWriter()

	err := git.Start()
	if err != nil {
		return errors.New("Failed to start git")
	}

	done := make(chan error, 1)
	go func() {
		done <- git.Wait()
	}()

	select {
	case <-time.After(240 * time.Second):
		if err := git.Process.Kill(); err != nil {
			log.Fatal("failed to kill: ", err)
		}
		log.Println("process killed as timeout reached")
	case err := <-done:
		if err != nil {
			log.Printf("process done with error = %v", err)
		} else {
			log.Print("process done gracefully without error")
		}
	}
	return err
}
