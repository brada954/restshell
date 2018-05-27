package shell

import (
	"testing"

	"github.com/pborman/getopt/v2"
)

type baseCommand struct {
	testInt    int
	testString string
}

func (b *baseCommand) Execute(args []string) error {
	return nil
}

func (b *baseCommand) AddOptions(set *getopt.Set) {
	return
}

type subCommand struct {
}

func (b *subCommand) Execute(args []string) error {
	return nil
}

func (b *subCommand) AddOptions(set *getopt.Set) {
}

func (b *subCommand) GetSubCommands() []string {
	var result = []string{}
	return result
}

func TestSubCommandDefinition(t *testing.T) {
	var subCmdType CommandWithSubcommands

	subCmdType = &subCommand{}
	_ = subCmdType.GetSubCommands()
}

func TestBaseCommandInterface(t *testing.T) {
	b := &baseCommand{}
	var cmd Command
	cmd = b
	if _, ok := cmd.(CommandWithSubcommands); ok {
		t.Errorf("baseCommand reported sub-commands in alternate mechanism")
	}
}

func TestSubCommandInterface(t *testing.T) {
	s := &subCommand{}
	var cmd Command
	cmd = s
	if sc, ok := cmd.(CommandWithSubcommands); !ok {
		t.Errorf("subCommand did not have sub-commands in alternate mechanism")
	} else {
		_ = sc.GetSubCommands()
	}
}

// Validate that an instance of command is separated from original Instance
// This test is designed to help future support for instances
func TestCmdInstance(t *testing.T) {
	name := "BASE"

	AddCommand(name, CategoryUtilities, &baseCommand{testInt: 111})

	if cmd, ok := cmdMap[name]; !ok {
		t.Errorf("Falied to look up command on first pass")
		return
	} else {
		if base, ok := cmd.(*baseCommand); ok {
			instance := *base
			instance.testInt = 123
		}
	}

	if cmd, ok := cmdMap[name]; !ok {
		t.Errorf("Falied to look up command on second pass")
	} else {
		if base, ok := cmd.(*baseCommand); ok {
			if base.testInt != 111 {
				t.Errorf("Instance modification affected original: %d!=%d", 111, base.testInt)
			}
		}
	}
}
