package shell

// Simple string map
type simpleMap struct {
	data map[string]string
}

func NewSimpleHistoryMap(m map[string]string) (HistoryMap, error) {
	return &simpleMap{
		data: m,
	}, nil
}

func (mh *simpleMap) GetNode(path string) (interface{}, error) {
	if v, ok := mh.data[path]; ok {
		return v, nil
	}
	return nil, ErrNotFound
}
