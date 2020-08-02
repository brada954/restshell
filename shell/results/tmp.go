package results

// func MakeNodeFromObject(o interface{}, err error) (INode, error) {
// 	if err != nil {
// 		return nil, err
// 	}

// 	if o == nil {
// 		return NewNodeNil(), nil
// 	}
// 	switch v := o.(type) {
// 	case string:
// 		return NewNodeString(v), nil
// 	case float64:
// 		return NewNodeFloat64
// 	case int32:
// 	case int64:
// 	}
// }

// func FirstNodeToValue(fn modifiers.ValueModifier, nodes[]INode) (node INode, err error) {
// 	node = nil
// 	err = nil

// 	if len(nodes) > 0 {
// 		result, err = fn(list)
// 	}
// 	return MakeNode(result, err)
// }

// func NodesToValue(fn modifiers.ValueModifier, nodes []INode) (node INode, err error) {
// 	list := make([]interface{}, 0)
// 	for _, n := range nodes {
// 		list = append(list, n.Value())
// 	}
// 	result, err := fn(list)
// 	return result, err
// }

/////////    JSON Result ///////////////

/////////    JSON Result ///////////////
