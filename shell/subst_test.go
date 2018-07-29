package shell

import (
	"fmt"
	"strings"
	"testing"
)

func TestSubstringMatch(t *testing.T) {
	fmt.Println("Beginning test...")
	result := SubstituteString("this %%newguid(3)%% is a %%newguid(1,short)%% test %%newguid(3)%%")
	fmt.Println("Result: ", result)
}

func TestToLowerSubstitution(t *testing.T) {
	result := SubstituteString("all text \"%%tolower(1,str,\"This Is A Test\")%%\" should be lower")
	if result != strings.ToLower(result) {
		t.Errorf("All text was not lower case: %s", result)
	}

	expected := "all text \"this is a test\" should be lower"
	if result != expected {
		t.Errorf("String missmatch: %s!=%s", expected, result)
	}
}

func TestToLowerSubstitutionWithInvalidValue(t *testing.T) {
	subString := "all text \"%%tolower(1,str,\"This Is #{}\")%%\" should be lower"
	result := SubstituteString(subString)

	// Regex failure should perform no substitution
	if result != subString {
		t.Errorf("String missmatch: %s!=%s", subString, result)
	}
}
