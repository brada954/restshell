package results

import "errors"

// TextDocument - Simple text document (one path "/")
type TextDocument struct {
	text string
}

// GetNodes -- return a collection of nodes
func (tr *TextDocument) GetNode(path string) (Node, error) {
	if path != "/" {
		return nil, errors.New("path not found")
	}
	return NewStringNode(tr.text), nil
}

// GetNodes -- return a collection of nodes
func (tr *TextDocument) GetNodes(path string) (NodeCollection, error) {
	if n,err := tr.GetNode(path) ; err != nil {
		return nil,err
	} else {
		return NodeCollection{n}, nil
	}
}