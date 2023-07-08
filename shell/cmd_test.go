package shell

import (
	"testing"
)

type baseCommand struct {
	testInt    int
	testString string
}

func (b *baseCommand) Execute(args []string) error {
	return nil
}

func (b *baseCommand) AddOptions(set CmdSet) {
}

type subCommand struct {
}

func (b *subCommand) Execute(args []string) error {
	return nil
}

func (b *subCommand) AddOptions(set CmdSet) {
}

func (b *subCommand) GetSubCommands() []string {
	var result = []string{}
	return result
}

func TestSubCommandInterfaceExtractedFromCommand(t *testing.T) {
	var cmd Command = &subCommand{}

	if subCmd, ok := cmd.(CommandWithSubcommands); ok {
		list := subCmd.GetSubCommands()
		if list == nil || len(list) != 0 {
			t.Errorf("sub command was expected to have zero length list of sub-commands")
		}
	} else {
		t.Errorf("command was expected to convert to CommandWithSubcommands but it did not")
	}
}

func TestBaseCommandInterfaceFailsTypeConvertionToCommandWithSubcommansd(t *testing.T) {
	var cmd Command = &baseCommand{}

	if _, ok := cmd.(CommandWithSubcommands); ok {
		t.Errorf("baseCommand was expected to not have sub-commands")
	}
}

// Validate that an instance of command is separated from original Instance
// This test is designed to help future support for instances
func TestCmdInstance(t *testing.T) {
	name := "BASE"

	AddCommand(name, CategoryUtilities, &baseCommand{testInt: 111, testString: "test"})

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
