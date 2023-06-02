# RESTResponse

A set of documents for different parts of a REST Response including body, headers, cookies, and statistics and connection metadata.

Generally, the body of a response may come in various forms which may be one of:
Json
Xml
Text
Opaque (HTML, binary, others) 

# Result package

The result package is an attempt to normalize different document types like xml and json to be usable in the generic assertion library. The idea is to have standard assert calls independent of document type versus assert langague for each document.

Long term thinking includes processing capabilities like iteration and sub-document processing through script snippets (think of sub routines).

JSON is the most flexible format for representing different forms of data like cookies, headers, and text, and since RestShell is primarily used with JSON scenarios, JSON will be a common format to represent non-body related data in a REST response.

The Opaque document type is a catchall for types that cannot be respresented in the other formats. Binary type can be replaced by future documents types to handle complex documents like jpg, png, pdf, office documents, etc. Document inspection and traversal features of these documents can be very interesting to have to examine metadata or document content. However, these documennt types are highly specialized and out of scope for the forseable future. The make great capabilities 3rd parties with expertise in these document formats to add.

The system designed should allow future expansion.

## Result Inspection

A result must support inspection to validate data, so a node traversal mechanism like XPath and JSONPath is fundamental mechanism to select data.
Assertions are based on providing a path to a value or object and an expectation about the element referenced by the path including negative assertions.

A common feature of data selection is that the result may be single result or list of results. A list of results will be difficult to distinguish from a single result that may be an array (i.e. list) when in native golang types. When a list of result is returned from a selection, the term collection will be used to distingush the result from a single result that may be an array.

The biggest challenge in designing the abstractions for selecting nodes is whether the caller should provide a hint that a collection is expected or a single result is expected. Providing the hint makes the implementation easier and less dependendent on implementations having to know the nuances of the path selection at the cost that implementation of the assertion logic will have to have assertions that are tailored to single results or not.

Initially, the design...

## Iterative Inspection

*Adanced concept not in immediate plans*

To operate on data within a document iteratively, the extraction of components of a document may need to be reconstructable to Result Documents for iterative processing or have a node abstraction supporting children.

# Node

A Node provides a native value representation for assertions as well as an abstraction to support further node traversal if not a leaf node.

A NodeCollection is a separate representation of a collection of nodes that can be iterated, typically this can be a set of nodes identified from a path filter.

A Node scaler is a leaf node that has a golang representation (json property, xml attribute or inner text).

A NodeArray is "leaf node" that has an array of native values or constructs (json array, no xml concept).

A NodeObject is a intermediate node that can reprent a colletion of properties (json object, xml element)


Nodes have a set of operations that can get a golang typed value represented by node or the node can be converted to a result document.

When a node is converted to a result document it may become a different document type than its original:

JSONDocument decomposes to: JsonNodeObject, TextNode, IntegerNode, FloatNumber, NullNode, JsonArray
XMLDocument decomposes to: XMLNodeObject, XMLCollection, StringNode (what about string collection?)

A Node can become a document based on the Node types override for building a document; by default it may become a json document unless the collection has a different implementation (such as xml collection may be returned by an xml document and such a collection can build a xml document)

===============================

(Question: can we introduce a document transformation at this level; what about xslt's for xml as well)

JsonNode, FloatNumber, NUllNode, IntegerNode -> JsonDocument
StringNode -> JsonDocument or TextDocument?
XMLNode -> XMLDocument


Converts to a Result

A result contains several result documents collected from an command:

--path
--path-{result type}

Result types need to be registered

type ResultContainer interface {
    GetDocument(doctype string) ResultDocument
}

type IResultDocument interface {
    GetNodes(string) []INode
}

type INode {
    ToResultDocument()
    Value()
}





Nodes:

xmlresult is a xml document

xmlnode
Holds an element, how does it generate recursive
element
properties() returns element child Nodes
ToResult() takes element and recreates a xml documents

xml document has nodes for element, array of any nodes, or scaler 

jsonresult is a set of possible types: map[string]interface{}, []interface{}, or scaler object
jsonnode
Holds a json object, json array, json scaler


Jsondocument has nodes for object, array any nodes, scaler

