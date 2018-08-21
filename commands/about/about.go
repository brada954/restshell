package about

import (
	"fmt"
	"strings"

	"github.com/brada954/restshell/shell"
)

type AboutCommand struct {
	// Place getopt option value pointers here
}

type TopicInterface interface {
	GetKey() string         // Key for lookup and sub-command
	GetTitle() string       // Title for help display
	GetDescription() string // Decription of key in lists
	GetAbout() string       // The text to display about the topic
}

var topicList []TopicInterface = []TopicInterface{
	NewAuthTopic(),
	NewBenchmarkTopic(),
	NewSubstitutionTopic(),
}

func NewAboutCommand() *AboutCommand {
	return &AboutCommand{}
}

func (cmd *AboutCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("topic")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *AboutCommand) Execute(args []string) error {
	// Validate arguments

	if len(args) == 0 {
		cmd.executeTopicList()
	} else {
		cmd.executeTopic(args[0])
	}
	return nil
}

func (cmd *AboutCommand) executeTopicList() {
	fmt.Fprintln(shell.ConsoleWriter(), "ABOUT {topic}")
	fmt.Fprintln(shell.ConsoleWriter(), "\nUse the ABOUT command to learn about the following topics:")
	fmt.Fprintln(shell.ConsoleWriter())
	for _, topic := range topicList {
		fmt.Fprintf(shell.ConsoleWriter(), "%s -- %s\n", topic.GetKey(), topic.GetDescription())
	}
	fmt.Fprintln(shell.ConsoleWriter(), "")
}

func (cmd *AboutCommand) executeTopic(key string) {
	for _, topic := range topicList {
		if strings.ToUpper(topic.GetKey()) != strings.ToUpper(key) {
			continue
		}

		fmt.Fprintf(shell.ConsoleWriter(), "%s (%s)\n\n", topic.GetTitle(), topic.GetKey())
		fmt.Fprintln(shell.ConsoleWriter(), topic.GetAbout())
		break
	}
}
