package shell

import (
	"testing"
)

func TestEscapeStringValueSuccess(t *testing.T) {
	expected := "This \\\\ test is \\\"everything\\\" there is!"
	input := "This \\ test is \"everything\" there is!"

	result := EscapeStringValue(input)
	if expected != result {
		t.Errorf("Expected (%s) but got (%s)", expected, result)
	}
}
