package about

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brada954/restshell/shell"
)

type AboutCommand struct {
	// Place getopt option value pointers here
}

func NewAboutCommand() *AboutCommand {
	return &AboutCommand{}
}

func (cmd *AboutCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("topic")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

// Execute -- Execute the command About command
func (cmd *AboutCommand) Execute(args []string) error {
	// Validate arguments
	if len(args) == 0 {
		return cmd.executeTopicList()
	} else {
		subTopic := ""
		if len(args) > 1 {
			subTopic = args[1]
		}
		return cmd.executeTopic(args[0], subTopic)
	}
}

func (cmd *AboutCommand) executeTopicList() error {
	fmt.Fprintln(shell.ConsoleWriter(), "ABOUT {topic}")
	fmt.Fprintln(shell.ConsoleWriter(), "\nUse the ABOUT command to learn about the following topics:")
	fmt.Fprintln(shell.ConsoleWriter())
	for _, topic := range shell.GetTopics() {
		fmt.Fprintf(shell.ConsoleWriter(), "%s -- %s\n", topic.GetKey(), topic.GetDescription())
	}
	fmt.Fprintln(shell.ConsoleWriter(), "")
	return nil
}

func (cmd *AboutCommand) executeTopic(key string, subTopic string) error {
	for _, topic := range shell.GetTopics() {
		if !strings.EqualFold(topic.GetKey(), key) {
			continue
		}

		if len(subTopic) > 0 {
			if st, ok := topic.(shell.SubTopicInterface); ok {
				return st.WriteSubTopic(shell.ConsoleWriter(), subTopic)
			} else {
				return errors.New("no sub-topics to display")
			}
		} else {
			fmt.Fprintf(shell.ConsoleWriter(), "%s (%s)\n\n", topic.GetTitle(), topic.GetKey())
			return topic.WriteAbout(shell.ConsoleWriter())
		}
	}
	return fmt.Errorf("%s topic was not found", key)
}
