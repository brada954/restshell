package shell

import "strings"

type textMap struct {
	data string
}

func NewTextHistoryMap(text string) (HistoryMap, error) {
	return &textMap{
		data: strings.TrimSpace(text),
	}, nil
}

func (tm *textMap) GetNode(path string) (interface{}, error) {
	if path == "/" {
		return tm.data, nil
	}
	return nil, ErrNotFound
}
