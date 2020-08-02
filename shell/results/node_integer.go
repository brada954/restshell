package results

import "fmt"

// IntegerNode --
type IntegerNode struct {
	value int64
}

// Value --
func (in *IntegerNode) Value() interface{} {
	return in.value
}

// String --
func (in *IntegerNode) String() string {
	return fmt.Sprintf("%v", in.value)
}

// NewInteger --
func NewInteger(val int64) *IntegerNode {
	return &IntegerNode{value: val}
}