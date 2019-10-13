package about

import (
	"fmt"
	"io"
)

type JsonPathTopic struct {
	Key         string
	Title       string
	Description string
	About       string
}

var localJsonPathTopic = &JsonPathTopic{
	Key:         "JSONPATH",
	Title:       "JSONPATH",
	Description: "JsonPath for accessing JSON elements",
	About: `JsonPath is a mechanism to access JSON objects within a JSON object. JsonPath is
used by the asserts and other commands to have a string paranter that can reference 
data from a JSON object.

JSON objects may be the results of a previous command stored in the history buffer it
may refer to a variable holding a JSON structure.

JsonPath can access a single property or potentially a collection of elements. However, 
many commands may expect a single property to be returned. Some asserts can get the 
length of a collection of elements for example. Otherwise, asserts have a --first
option to get the first element if the assert requires a property.

Please refer to the following links about the standard and some examples of uses:
https://goessner.net/articles/JsonPath/
https://www.baeldung.com/guide-to-jayway-jsonpath
`,
}

func NewJsonPathTopic() *JsonPathTopic {
	return localJsonPathTopic
}

func (a *JsonPathTopic) GetKey() string {
	return a.Key
}

func (a *JsonPathTopic) GetTitle() string {
	return a.Title
}

func (a *JsonPathTopic) GetDescription() string {
	return a.Description
}

func (a *JsonPathTopic) WriteAbout(o io.Writer) error {
	fmt.Fprintf(o, a.About)
	return nil
}
