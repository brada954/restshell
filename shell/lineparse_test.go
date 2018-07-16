package shell

import (
	"errors"
	"testing"
)

func TestParseQuotes(t *testing.T) {
	testAString(t, "123 abc", []string{"123", "abc"})
	testAString(t, "\"Test In Quotes\"", []string{"Test In Quotes"})
	testAString(t, "123\"abc\"", []string{"123abc"})
	testAString(t, "123\"abc", []string{"123abc"})
	testAString(t, "123\\\"abc", []string{"123\\abc"})
	testAString(t, "123 \"\\\"abc\\\"", []string{"123", "\"abc\""})
	testAString(t, "-url=\"http://this web site\"",
		[]string{"-url=http://this web site"})
	testAString(t, "-url=http://this web site",
		[]string{"-url=http://this", "web", "site"})
	testAString(t, "123 \"abc\"", []string{"123", "abc"})
	testAString(t, "123 \"\\\"abc\\\"\"", []string{"123", "\"abc\""})
	testAString(t, "123 \" abc \"\"", []string{"123", " abc "})
	testAString(t, "buy   -v -123", []string{"buy", "-v", "-123"})
}

func testAString(t *testing.T, line string, expected []string) {
	args := LineParse(line)
	if len(args) != len(expected) {
		t.Errorf("Invalid number of tokens %d!=%d", len(args), 2)
	}

	for i := 0; i < min(len(args), len(expected)); i++ {
		if args[i] != expected[i] {
			t.Errorf("arg[%d] does not match: %s!=%s", i, args[i], expected[i])
		}
	}
}

func TestVariableSubstitution(t *testing.T) {
	initGlobalStore()
	SetGlobal("xyz", "zyx")
	SetGlobal("abc", errors.New("This is a string"))
	SetGlobal("nil", nil)
	SetGlobal("friend", " friend ")

	testSubstitution(t, "xyz abc 123", "xyz abc 123")
	testSubstitution(t, "xyz %%abc%% 123", "xyz %%abc%% 123")
	testSubstitution(t, "%%xyz%% abc 123", "zyx abc 123")
	testSubstitution(t, "this is my%%friend%%here", "this is my friend here")
	testSubstitution(t, "this is a test of %%nil%%", "this is a test of %%nil%%")
	testSubstitution(t, "this test of \"%%xyz%% is tired\"", "this test of \"zyx is tired\"")
}

func TestCompleteVariableSubstitution(t *testing.T) {
	initGlobalStore()

	verifyCompleteSubstitution(t, "abcd efg hijk lmknop", true)
	verifyCompleteSubstitution(t, "abcd efg %%hijk%% lmknop", false)
	verifyCompleteSubstitution(t, "abcd %efg% hijk lmknop", true)
	verifyCompleteSubstitution(t, "abcd %efghijk lm%%$kn%%op", false)
}

func testSubstitution(t *testing.T, input string, expected string) {
	result := PerformVariableSubstitution(input)
	if expected != result {
		t.Errorf("Failed substition: %s != %s", expected, result)
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func verifyCompleteSubstitution(t *testing.T, input string, complete bool) {
	if IsVariableSubstitutionComplete(input) {
		if !complete {
			t.Errorf("Variables to be substuted exist unexpectedly: %s", input)
		}
	} else {
		if complete {
			t.Errorf("Variables were not found as expected: %s", input)
		}
	}
}
