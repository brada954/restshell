package functions

import (
	"strings"
	"testing"

	"github.com/brada954/restshell/shell"
)

func TestToLowerSubstitution(t *testing.T) {
	result := shell.PerformVariableSubstitution("all text \"%%tolower(1,str,\"This Is A Test\")%%\" should be lower")
	if result != strings.ToLower(result) {
		t.Errorf("All text was not lower case: %s", result)
	}

	expected := "all text \"this is a test\" should be lower"
	if result != expected {
		t.Errorf("String missmatch: %s!=%s", expected, result)
	}
}

func TestToLowerSubstitutionWithInvalidValue(t *testing.T) {
	subString := "all text \\\"%%tolower(1,str,\"This Is #{}\")%%\\\" should be lower"
	result := shell.PerformVariableSubstitution(subString)

	// Regex failure should perform no substitution
	if result != subString {
		t.Errorf("String missmatch: %s!=%s", subString, result)
	}
}
