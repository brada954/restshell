package results

import (
	"testing"
)

func TestGetNodeOnSimpleTextDocument(t *testing.T) {
	var expected = "This is a test document"
	var doc = &TextDocument{text : expected}
	
	n,err := doc.GetNode("/")
	{
		if err != nil {
			t.Errorf("Unexpected failure of GetNode(%s) : %s", "/", err.Error())
		}

		if n.Value() != expected {
			t.Errorf("Expected string (%s) does not match node value: (%s)", expected, n.Value())
		}
	}
}


func TestGetNodesOnSimpleTextDocument(t *testing.T) {
	var expected = "This is a test document"
	var doc = &TextDocument{text : expected}
	
	nc,err := doc.GetNodes("/")
	{
		if err != nil {
			t.Errorf("Unexpected failure of GetNode(%s) : %s", "/", err.Error())
		}

		if len(nc) != 1 {
			t.Errorf("Expected one node, but %d were returned", len(nc))
			return
		}

		n := nc[0]
		if n.Value() != expected {
			t.Errorf("Expected string (%s) does not match node value: (%s)", expected, n.Value())
		}
	}
}

func TestBadPathGetNodeOnSimpleTextDocument(t *testing.T) {
	var expected = "This is a test document"
	var expectedError = "path not found"
	var doc = &TextDocument{text : expected}
	
	_,err := doc.GetNode("")
	{
		if err == nil {
			t.Errorf("Unexpected success of GetNode(%s)", "")
		}

		if err.Error() != expectedError {
			t.Errorf("Expected error (%s) does not match error: (%s)", expectedError, err)
		}
	}
}