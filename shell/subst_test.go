package shell

import (
	"fmt"
	"testing"
)

func TestSubstringMatch(t *testing.T) {
	fmt.Println("Beginning test...")
	result := SubstituteString("this %%newguid(3)%% is a %%newguid(1,short)%% test %%newguid(3)%%")
	fmt.Println("Result: ", result)
}
