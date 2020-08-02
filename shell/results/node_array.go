package results

// ArrayNode - a container of an array of scalers
type ArrayNode struct {
	value []interface{}
}

// Value -- golang value of an array : []interface{}
func (fn *ArrayNode) Value() interface{} {
	return fn.value
}

func (fn *ArrayNode) Length() int {
	return len(fn.value)
}

func NewArrayNode(values []interface{}) *ArrayNode {
	return &ArrayNode{value: values}
}
