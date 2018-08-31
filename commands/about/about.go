package about

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/brada954/restshell/shell"
)

type AboutCommand struct {
	// Place getopt option value pointers here
}

// TopicInterface -- THe minumum supported interface for about topics
type TopicInterface interface {
	GetKey() string             // Key for lookup and sub-command
	GetTitle() string           // Title for help display
	GetDescription() string     // Decription of key in lists
	WriteAbout(io.Writer) error // The text to display about the topic
}

// SubTopicInterface -- Some about topics may have sub topics
type SubTopicInterface interface {
	WriteSubTopic(io.Writer, string) error
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

// Execute --
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
	return nil
}

func (cmd *AboutCommand) executeTopicList() error {
	fmt.Fprintln(shell.ConsoleWriter(), "ABOUT {topic}")
	fmt.Fprintln(shell.ConsoleWriter(), "\nUse the ABOUT command to learn about the following topics:")
	fmt.Fprintln(shell.ConsoleWriter())
	for _, topic := range topicList {
		fmt.Fprintf(shell.ConsoleWriter(), "%s -- %s\n", topic.GetKey(), topic.GetDescription())
	}
	fmt.Fprintln(shell.ConsoleWriter(), "")
	return nil
}

func (cmd *AboutCommand) executeTopic(key string, subTopic string) error {
	for _, topic := range topicList {
		if strings.ToUpper(topic.GetKey()) != strings.ToUpper(key) {
			continue
		}

		if len(subTopic) > 0 {
			if st, ok := topic.(SubTopicInterface); ok {
				if err := st.WriteSubTopic(shell.ConsoleWriter(), subTopic); err != nil {
					return err
				}
			} else {
				return errors.New("No sub-topics to display")
			}
		} else {
			fmt.Fprintf(shell.ConsoleWriter(), "%s (%s)\n\n", topic.GetTitle(), topic.GetKey())
			return topic.WriteAbout(shell.ConsoleWriter())
		}
		break
	}
	return nil
}
