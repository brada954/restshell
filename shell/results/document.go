package results

// ResultDocumentType -- the type of document
type ResultDocumentType int

// ResultDocument types
const (
	JSONResultDocument = iota + 1
	XMLResultDocument
	TextResultDocument
	KeyValueResultDocument
)

// ResultDocument -- interface for result documents
type ResultDocument interface {

	// GetNode -- Given a path return a single node
	// returns 
	//   a Node
	//   an error if zero or more than one node is found
	GetNode(path string) (Node, error)
	
	// GetNodes -- Given a path return the identified node(s)
	// returns
	//   a collection of nodes
	//   an error code
	GetNodes(path string) (NodeCollection, error)
}

// DocumentMaker -- interface to create a document from a node
type DocumentMaker interface {
	ToDocument(node Node) error
}

