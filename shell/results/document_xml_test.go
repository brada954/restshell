package results

import (
	"strings"
	"testing"
)

func TestCanXMLDocumentGetNodeReturnNodeStringForPath(t *testing.T) {
	xml := "<root><header>Title</header><body>This is a test document</body><list><listitem>Item1</listitem><listitem>item2</listitem></list></root>"

	doc, err := MakeXmlResult(strings.NewReader(xml))
	if err != nil {
		t.Errorf("Unable to setup XML document")
		return
	}

	n, err := doc.GetNode("/root/header")

	if err != nil {
		t.Errorf("Expected success but got the error: %s", err.Error())
	}
	if n.Value() != "Title" {
		t.Errorf("Expected title for header but got: %v", n.Value())
	}
}
