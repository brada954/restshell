package kubectl

import (
	"errors"
	"log"
	"os/exec"
	"time"

	"github.com/brada954/restshell/shell"
)

var ()

type KubectlCommand struct {
	// Globals

	// Options
}

func NewKubectlCommand() *KubectlCommand {
	return &KubectlCommand{}
}

func (cmd *KubectlCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("kubectl_command [options]")
	set.SetUsage(func() {
		set.PrintUsage(shell.ConsoleWriter())
	})
	shell.AddCommonCmdOptions(set, shell.CmdDebug)
}

func (cmd *KubectlCommand) Execute(args []string) error {

	kube := exec.Command("kubectl", args...)

	kube.Stdout = shell.OutputWriter()
	kube.Stderr = shell.ErrorWriter()

	err := kube.Start()
	if err != nil {
		return errors.New("Failed to start kubectl")
	}

	done := make(chan error, 1)
	go func() {
		done <- kube.Wait()
	}()

	select {
	case <-time.After(120 * time.Second):
		if err := kube.Process.Kill(); err != nil {
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
