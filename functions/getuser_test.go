package functions

import (
	"testing"
)

func TestGetUserData(t *testing.T) {
	c1, err := getRandomUserData()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(c1.Name.First) == 0 {
		t.Errorf("Name was not returned from getUserData")
	}

	c2, err := getRandomUserData()
	if err != nil {
		t.Errorf("Unexpected error on call2: %s", err)
	}
	if len(c2.Name.First) == 0 {
		t.Errorf("Name was not returned from second getUserData")
	}
}
