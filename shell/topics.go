package shell

import (
	"io"
)

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

var topicList []TopicInterface = make([]TopicInterface, 0, 10)

func GetTopics() []TopicInterface {
	return topicList
}

func AddAboutTopic(topic TopicInterface) {
	if topicsContains(topic.GetKey()) {
		panic("Adding a duplicate about topic is not allowed")
	}
	topicList = append(topicList, topic)

}

func topicsContains(key string) bool {
	for _, t := range topicList {
		if t.GetKey() == key {
			return true
		}
	}
	return false
}
