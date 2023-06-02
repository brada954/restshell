package results

// ResultDocument -- interface for inspecting document elements using a path
// to retrieve nodes used by assertions and scripting components
type ResultDocument interface {

	// GetNode -- Given a path return a single node
	// returns
	//   a Node
	//   an error if zero or more than one node is found
	GetNode(path string) (Node, error)

	// GetNodes -- Given a path return the matching node(s)
	// returns
	//   a collection of nodes
	//   an error code
	GetNodes(path string) (NodeCollection, error)
}

type DocumentDeconstructor interface {
	// MakeDocument creates a document from the node(s) identified by the path
	// When multiple nodes are identified by the path, the root document represents
	// an array of the results
	MakeDocument(path string) (ResultDocument, error)
}
