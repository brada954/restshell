package modifiers

import (
	"testing"
)

func TestRegexModifier(t *testing.T) {
	testRegexMod(t, "failed", "This is a error message for a failed operation", "failed")
	testRegexMod(t, "failed", "This is an error message", "")
	testRegexMod(t, "\\((\\d*)\\)$", "Error code (27)", "27")
	testRegexMod(t, "test (this( app)) end", "my test this app end", "this app app")
	testRegexMod(t, "test (this( app)) end", "my texst this app end", "")
}

func testRegexMod(t *testing.T, pattern string, value interface{}, expected string) {
	valueModifierFunc := NullModifier
	valueModifierFunc = MakeRegExModifier(pattern, valueModifierFunc)

	newval, err := valueModifierFunc(value)
	if err != nil {
		t.Errorf("Regex modifier error %s!=%s: %s", expected, newval, err.Error())
	}
	if newval != expected {
		t.Errorf("Regex Modifier Failed: %s!=%s", expected, newval)
	}
}
