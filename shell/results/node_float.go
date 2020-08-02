package results

import "fmt"

type FloatNode struct {
	value float64
}

func (fn *FloatNode) Value() interface{} {
	return fn.value
}

func (fn *FloatNode) String() string {
	return fmt.Sprintf("%v", fn.value)
}

func NewFloat(val float64) *FloatNode {
	return &FloatNode{value: val}
}