package shell

import (
	"fmt"
	"testing"
)

func TestColumnizeZero(t *testing.T) {
	result := ColumnizeTokens([]string{}, 3, 10)
	if len(result) != 0 {
		t.Errorf("Unexpected lines returned; expected empty list")
	}
}

func TestColumnizeOne(t *testing.T) {
	result := ColumnizeTokens([]string{"test"}, 3, 10)
	if len(result) != 1 {
		t.Errorf("Unexpected number lines returned; 1 != %d", len(result))
	}

	expected := "test      "
	if len(result) == 1 && result[0] != expected {
		t.Errorf("Line does not match expectation: (%s)!=(%s)", expected, result[0])
	}
}

func TestColumnizeTwo(t *testing.T) {
	result := ColumnizeTokens([]string{"test", "apple"}, 3, 10)
	if len(result) != 1 {
		t.Errorf("Unexpected number lines returned; 1 != %d", len(result))
	}

	expected := "test      apple     "
	if len(result) == 1 && result[0] != expected {
		t.Errorf("Line does not match expectation: (%s)!=(%s)", expected, result[0])
	}
}

func TestColumnizeThree(t *testing.T) {
	result := ColumnizeTokens([]string{"test", "apple", "last"}, 2, 10)
	if len(result) != 2 {
		t.Errorf("Unexpected number lines returned; 2 != %d", len(result))
	}

	expected := "test      apple     "
	if len(result) > 0 && result[0] != expected {
		t.Errorf("Line does not match expectation: (%s)!=(%s)", expected, result[0])
	}

	expected = "last      "
	if len(result) >= 2 && result[1] != expected {
		t.Errorf("Line does not match expectation: (%s)!=(%s)", expected, result[1])
	}
}

func TestColumnizeFive(t *testing.T) {
	result := ColumnizeTokens([]string{"test", "apple", "second", "line", "last"}, 2, 10)
	if len(result) != 3 {
		t.Errorf("Unexpected number lines returned; 3 != %d", len(result))
	}

	expected := "test      apple     "
	if len(result) > 0 && result[0] != expected {
		t.Errorf("Line 1 does not match expectation: (%s)!=(%s)", expected, result[0])
	}

	expected = "second    line      "
	if len(result) > 1 && result[1] != expected {
		t.Errorf("Line 2 does not match expectation: (%s)!=(%s)", expected, result[1])
	}

	expected = "last      "
	if len(result) > 2 && result[2] != expected {
		t.Errorf("Line 3 does not match expectation: (%s)!=(%s)", expected, result[2])
	}
}

func TestIsStringBinaryWithBinary(t *testing.T) {
	var text = "{ \013this\013 is a \003test\000 that\013\013\030 needs 10\001 binary ch\013arac\013ters\034}"
	text = text + "extra c\030racters\030 to get this over 100 be\001ca\001use there"
	if !isStringBinary(text) {
		t.Errorf("Unexpected failure to validate binary string as binary")
		fmt.Println(text)
	}
}

func TestIsStringBinaryWithShortText(t *testing.T) {
	var text = "{ this is a \003test that needs 10\001 binary characters}"
	if isStringBinary(text) {
		t.Errorf("Unexpected failure to validate text string as not binary")
		fmt.Println(text)
	}
}

func TestIsStringBinaryWithLongText(t *testing.T) {
	var text = "{ this is a \003test that needs 10\001 binary chara\001cters}"
	text = text + " some extra characters to get this over 100; 123123"
	if isStringBinary(text) {
		t.Errorf("Unexpected failure to validate text string as not binary")
		fmt.Println(text)
	}
}

func outputStrings(lines []string) {
	for _, s := range lines {
		fmt.Printf("Line: (%s)\n", s)
	}
	fmt.Println("")
}
