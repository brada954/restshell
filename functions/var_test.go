package functions

import (
	"testing"

	"github.com/brada954/restshell/shell"
)

func TestQuoteSubstitution(t *testing.T) {
	testQuoteSubstitution(t, "{ }", "\"{ }\"")
	testQuoteSubstitution(t, "{\"test\" : 123, \"test2\" : \"string\"}", "\"{\\\"test\\\" : 123, \\\"test2\\\" : \\\"string\\\"}\"")
}

func testQuoteSubstitution(t *testing.T, input string, expected string) {
	shell.SetGlobal("testvar", input)
	result := shell.PerformVariableSubstitution("JSON: %%quote(1,testvar)%%.")
	expected = "JSON: " + expected + "."
	if result != expected {
		t.Errorf("Text was not expected: %s!=%s", expected, result)
	}
}
