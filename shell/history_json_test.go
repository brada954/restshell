package shell

import (
	"strings"
	"testing"
)

func TestCanJsonGetNodeParseJsonString(t *testing.T) {
	var parseString = "\"The blue dog\""
	if hmap, err := NewJsonHistoryMap(parseString); err != nil {
		t.Errorf("NewJsonHistoryMap failed to parse \"%s\" and returned error: %s", parseString, err.Error())
	} else {
		if n, err2 := hmap.GetNode("$"); err != nil {
			t.Errorf("GetNode failed to extract string for $ path and returned error: %s", err2.Error())
		} else if n != strings.Trim(parseString, "\"") {
			t.Errorf("Expected to parse out \"%s\" but got \"%s\"", parseString, n)
		}
	}
}
