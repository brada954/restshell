////////////////////////////////////////////////////////////////////////////////
// Nodes
//   Nodes encapsulate the values that can be extracted from result documents
//   and used by the various components of RestShell to process data.
//
//   All nodes have a value respresentation which is a natural golang
//   construct (string, integer, interface, slice, etc)
//
//   Nodes can also represents other data elements which may include collections,
//   arrays and objects.
////////////////////////////////////////////////////////////////////////////////
package results

import (
	"errors"
)

// Node -- abstract representation of a xpath/jsonpath like node that can
// provide a golang type representing the data or can generate a new ResultDocument
type Node interface {
	// Value - return the golang data type
	Value() interface{}
}

// NodeCollection -- A collection of Nodes as scanned from a document
type NodeCollection []Node

// NodeScaler -- represents a scaler value and it can be represented as a string
type NodeScaler interface {
	Node
	String() string
}

// NodeArray is a specfic Node in the document representing an array of values, objects or
// arrays. The NodeArray just has a length.
type NodeArray interface {
	Node
	Length()
}

// NodeObject represents an object containing a set of properties
type NodeObject interface {
	Node
	Properties() []string
}

// ToNode -- convert GetNodes result to a single element
func (collection NodeCollection) ToNode() (Node, error) {
	if collection == nil || len(collection) == 0 {
		return nil, errors.New("Node not found")
	}

	if len(collection) > 1 {
		return nil, errors.New("too many nodes to for a single node operation")
	}
	return collection[0], nil
}
