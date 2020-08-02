package results

// StringNode implements NodeScaler for string value
type StringNode struct {
	value string
}

// Value -- returns the string value
func (sn *StringNode) Value() interface{} {
	return sn.value
}

// String -- returns the string respresentation of value
func (sn *StringNode) String() string {
	return sn.value
}

// NewStringNode returns a Node for a string
func NewStringNode(txt string) *StringNode {
	return &StringNode{value: txt}
}