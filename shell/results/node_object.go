package results

// ObjectNode - an node representing an object or xml element
type ObjectNode struct {
	value      interface{}
	properties []string
}

func (on *ObjectNode) Value() interface{} {
	return on.value
}

func (on *ObjectNode) Properties() []string {
	return on.properties
}

func NewObjectNode(obj interface{}, props []string) *ObjectNode {
	return &ObjectNode{value: obj, properties: props}
}