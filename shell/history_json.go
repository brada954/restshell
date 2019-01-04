package shell

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Json map
type jsonMap struct {
	data interface{}
}

// MakeJsonHistoryMap -- Create a HistoryMap for json content
func NewJsonHistoryMap(data string) (HistoryMap, error) {

	resultMap, err := makeHistoryMapFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &jsonMap{
		data: resultMap,
	}, nil
}

// GetNode -- Get a JSON node from a map structure mapped from a json object or array
func (jm *jsonMap) GetNode(path string) (interface{}, error) {
	return getNodeImpl(path, jm.data)
}

func makeHistoryMapFromJSON(data string) (interface{}, error) {
	var f interface{}

	err := json.Unmarshal([]byte(data), &f)
	if err != nil {
		return nil, err
	}

	if m, ok := f.(map[string]interface{}); ok {
		return m, nil
	}

	if m, ok := f.([]interface{}); ok {
		return m, nil
	}

	m := make(map[string]interface{}, 1)
	m["/"] = f
	return m, nil
}

// getNodeImpl -- Recursive parser of json node
func getNodeImpl(path string, i interface{}) (interface{}, error) {
	if path == "/" {
		return i, nil
	}

	parts := strings.SplitN(path, ".", 2)
	if len(parts) <= 0 {
		return nil, ErrInvalidPath
	}

	arrIndex := -1
	arrParts := strings.SplitN(parts[0], "[", 2)
	if len(arrParts) > 1 {
		index, err := strconv.Atoi(strings.Trim(arrParts[1], "]"))
		if err != nil {
			return nil, ErrInvalidPath
		}
		arrIndex = index
	}

	if arrIndex >= 0 {
		var data interface{}

		if len(arrParts[0]) > 0 {
			m, ok := i.(map[string]interface{})
			if !ok {
				return nil, ErrUnexpectedType
			}
			data = m[arrParts[0]]
		} else {
			if a, ok := i.([]interface{}); !ok {
				return nil, ErrUnexpectedType
			} else {
				data = a
			}
		}

		if data == nil {
			return nil, ErrNotFound
		}

		if len(parts) == 1 {
			switch t := data.(type) {
			case []string:
				if arrIndex < len(t) {
					return t[arrIndex], nil
				}
				return nil, ErrArrayOutOfBounds
			case []float64:
				if arrIndex < len(t) {
					return t[arrIndex], nil
				}
				return nil, ErrArrayOutOfBounds
			case []interface{}:
				if arrIndex < len(t) {
					dv := t[arrIndex]
					return dv, nil
				}
				return nil, ErrArrayOutOfBounds
			}
			return nil, ErrUnexpectedType
		} else {
			switch t := data.(type) {
			case []map[string]interface{}:
				if arrIndex < len(t) {
					return getNodeImpl(parts[1], t[arrIndex])
				}
				return nil, ErrArrayOutOfBounds
			case []interface{}:
				if arrIndex < len(t) {
					return getNodeImpl(parts[1], t[arrIndex])
				}
				return nil, ErrArrayOutOfBounds
			}
			return nil, ErrUnexpectedType
		}
	} else {
		m, ok := i.(map[string]interface{})
		if !ok {
			return nil, ErrUnexpectedType
		}

		if len(parts) > 1 {
			data := m[parts[0]]
			if data == nil {
				return nil, ErrNotFound
			}
			switch t := m[parts[0]].(type) {
			case map[string]interface{}:
				return getNodeImpl(parts[1], t)
			}
			return nil, ErrUnexpectedType
		} else {
			data, ok := m[parts[0]]
			if !ok && parts[0] == "/" {
				data = m
			}
			if data == nil {
				if !ok {
					return nil, ErrNotFound
				}
				return nil, nil
			}
			switch t := data.(type) {
			case string:
				return t, nil
			case int:
				return t, nil
			case float64:
				return t, nil
			default:
				return data, nil
			}
		}
	}
}
