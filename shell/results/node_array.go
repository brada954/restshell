package results

// ArrayNode - a container of an array of scalers
type ArrayNode struct {
	length int
}

// Value -- always nil
func (fn *ArrayNode) Value() interface{} {
	return nil
}

// Length - the length of the array
func (fn *ArrayNode) Length() int {
	return fn.length
}

// NewArrayNode -- create a Node for an array
func NewArrayNode(length int) *ArrayNode {
	return &ArrayNode{length: length}
}
