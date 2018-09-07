package functions

import (
	"fmt"
	"testing"

	"github.com/brada954/restshell/shell"
)

func TestSubstringMatch(t *testing.T) {
	fmt.Println("Beginning test...")
	result := shell.PerformVariableSubstitution("this %%newguid(3)%% is a %%newguid(1,short)%% test %%newguid(3)%%")
	fmt.Println("Result: ", result)
}
