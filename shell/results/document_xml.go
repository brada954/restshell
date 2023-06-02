package results

import (
	"errors"
	"io"

	"github.com/antchfx/xmlquery"
)

// XmlResult -- XML implementation of a result
type XmlResult struct {
	data *xmlquery.Node
}

// MakeXmlResult -- Given a reader return a XmlResult
func MakeXmlResult(r io.Reader) (*XmlResult, error) {
	if doc, err := xmlquery.Parse(r); err == nil {
		return &XmlResult{
			data: doc,
		}, nil
	} else {
		return nil, err
	}
}

func (x *XmlResult) GetNode(path string) (Node, error) {
	if collection, err := x.GetNodes(path); err != nil {
		return nil, err
	} else if len(collection) == 1 {
		return collection[0], nil
	} else {
		return nil, errors.New("Invalid number of nodes returned")
	}
}

// GetNodes -- return Nodes for the XmlResult
func (x *XmlResult) GetNodes(path string) (nodes NodeCollection, rtnerr error) {

	nodes = make([]Node, 0)
	rtnerr = nil

	defer func() {
		if r := recover(); r == nil {
			return // Pass-thru existing error code
		} else {
			nodes = make([]Node, 0)
			rtnerr = errors.New("Error with XPATH: " + path)
		}
	}()

	if x.data == nil {
		return nil, errors.New("not found") //ErrNotFound
	}

	nodeList, err := xmlquery.QueryAll(x.data, path)
	array := make([]Node, 0)
	if err != nil {
		return nil, err
	}
	if len(nodeList) >= 1 {
		for _, v := range nodeList {
			array = append(array, NewStringNode(v.InnerText()))
		}
		return array, nil
	} else {
		return nil, errors.New("no elements Found")
	}
}

// XmlNode --
type XmlNode struct {
	xmlquery.Node
	properties []string
}

// func (xn *XmlNode) ToResult() IResultDocument {
// 	return nil
// }

func (xn *XmlNode) Value() interface{} {

	switch xn.Node.Type {
	case xmlquery.DocumentNode:
	case xmlquery.DeclarationNode: // Not used
	case xmlquery.ElementNode:
		//return xn.newNodeObject()
	case xmlquery.TextNode:
	case xmlquery.CommentNode: // Not used
	case xmlquery.AttributeNode:
	}
	return xn // We do not have a scaler value
}

// func (xn *XmlNode) Length() int {
// 	return 0
// }

func (xn *XmlNode) Properties() []string {
	return xn.properties
}

// newNodeObject -- node object is an array of child element names
// func (xn *XmlNode) newNodeObject() NodeObject {
// 	result := make([]string, 0)
// 	for elem := xn.Node.FirstChild; elem != nil; elem = elem.NextSibling {
// 		result = append(result, elem.Data)
// 	}
// 	return &XmlNode{properties: result}
// }
