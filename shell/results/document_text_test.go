package results

import (
	"reflect"
	"testing"
)

func TestCanTextDocumentGetNodeOfRootReturnNodeString(t *testing.T) {
	expected := "This is a test document"
	doc := &TextDocument{text: expected}

	n, err := doc.GetNode("/")

	validateGetNodeSuccess(t, n, err, expected)
}

func TestCanTextDocumentGetNodeReturnErrorForEmptyPath(t *testing.T) {
	expected := "This is a test document"
	doc := &TextDocument{text: expected}

	n, err := doc.GetNode("")

	validateGetNodeError(t, n, err, "path not found")
}

func TestCanTextDocumentGetNodeReturnErrorForInvalidPath(t *testing.T) {
	expected := "This is a test document"
	doc := &TextDocument{text: expected}

	n, err := doc.GetNode("xyz")

	validateGetNodeError(t, n, err, "path not found")
}

func TestCanTextDocumentGetNodesOfRootReturnCollectionOfNodeString(t *testing.T) {
	expected := "This is a test document"
	doc := &TextDocument{text: expected}

	n, err := doc.GetNodes("/")

	validateGetNodeCollectionSuccess(t, n, err, []string{expected})
}

func TestTextDocumentWhenGettingNodesForEmptyPathThenErrorReturned(t *testing.T) {
	var doc = &TextDocument{text: "This is a test document"}

	n, err := doc.GetNodes("")

	validateGetNodeCollectionError(t, n, err, "path not found")
}

func validateGetNodeSuccess(t *testing.T, n Node, err error, expected string) {
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if n == nil {
		t.Errorf("Expected text (%s) for node content but nil was returned for node", expected)
	}
	if n.Value() != expected {
		t.Errorf("Expected text (%s) does not match returned text: (%s)", expected, n.Value())
	}
}

func validateGetNodeError(t *testing.T, n Node, err error, expected string) {
	if err == nil {
		t.Errorf("Unexpected success when expecting: %s", expected)
	}
	if err.Error() != expected {
		t.Errorf("Expected error (%s) but got: (%s)", expected, err)
	}
	if n != nil {
		t.Errorf("Expected a nil node but received: %v", n.Value())
	}
}

func validateGetNodeCollectionSuccess(t *testing.T, collection NodeCollection, err error, expected []string) {
	if err != nil {
		t.Errorf("Unexpected failure: %s", err.Error())
	}

	if len(collection) != len(expected) {
		t.Errorf("Expected collecction of (%d) node(s), but have %d", len(expected), len(collection))
	}

	for idx, n := range collection {
		if nt, ok := n.(*StringNode); !ok {
			t.Errorf("Unexpected node type: %s", reflect.TypeOf(nt))
		}
		if n.Value() != expected[idx] {
			t.Errorf("Expected text (%s) does not match returned text: (%s)", expected[idx], n.Value())
		}
	}
}

func validateGetNodeCollectionError(t *testing.T, collection NodeCollection, err error, expected string) {
	if err == nil {
		t.Errorf("Unexpected success when expecting: %s", expected)
	}
	if err.Error() != expected {
		t.Errorf("Expected error (%s) but got different error: (%s)", expected, err.Error())
	}
	if collection != nil {
		t.Errorf("Expected a nil collection but received collection of length: %d", len(collection))
	}
}
