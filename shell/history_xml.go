package shell

import (
	"errors"
	"strings"

	"github.com/antchfx/xmlquery"
)

// XML map
type xmlMap struct {
	data *xmlquery.Node
}

func NewXmlHistoryMap(data string) (HistoryMap, error) {
	data = strings.TrimSpace(data)
	if doc, err := xmlquery.Parse(strings.NewReader(data)); err == nil {
		return &xmlMap{
			data: doc,
		}, nil
	} else {
		return nil, err
	}
}

// GetNode - given an xpath return the node or nodes returned with
// the inner text
func (xm *xmlMap) GetNode(path string) (result interface{}, rtnerror error) {

	defer func() {
		if r := recover(); r == nil {
			return // Pass-thru existing error code
		} else {
			rtnerror = errors.New("Error with XPATH: " + path)
		}
	}()

	if xm.data == nil {
		return nil, ErrNotFound
	}

	nodeList, err := xmlquery.QueryAll(xm.data, path)
	if err != nil {
		return nil, err
	}
	if len(nodeList) > 1 {
		array := make([]string, 0)
		for _, v := range nodeList {
			array = append(array, v.InnerText())
		}
		return array, nil
	} else if len(nodeList) == 1 {
		return nodeList[0].InnerText(), nil
	} else {
		return nil, ErrNotFound
	}
}
