package shell

import (
	"fmt"
	"testing"
)

func TestSubstringMatch(t *testing.T) {
	fmt.Println("Beginning test...")
	// Need a real test; this actually doesn't substitute as the functions package is not included
	result := PerformVariableSubstitution("this %%newguid(3)%% is a %%newguid(1,short)%% test %%newguid(3)%%")
	fmt.Println("Result: ", result)
}
