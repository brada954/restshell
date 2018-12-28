package shell

import (
	"errors"
	"strings"

	"github.com/subchen/go-xmldom"
)

// XML map
type xmlMap struct {
	data *xmldom.Document
}

func NewXmlHistoryMap(data string) (HistoryMap, error) {
	wrapper := xmldom.Must(xmldom.ParseXML("<assertwrapper></assertwrapper>"))
	data = strings.TrimSpace(data)
	if doc, err := xmldom.ParseXML(data); err == nil {
		wrapper.Root.AppendChild(doc.Root)
		return &xmlMap{
			data: wrapper,
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

	root := xm.data.Root
	nodeList := root.Query(path)
	if len(nodeList) > 1 {
		array := make([]string, 0)
		for _, v := range nodeList {
			array = append(array, v.Text)
		}
		return array, nil
	} else if len(nodeList) == 1 {
		return nodeList[0].Text, nil
	} else {
		return nil, ErrNotFound
	}
}
