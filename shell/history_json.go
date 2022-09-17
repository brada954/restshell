package shell

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

// JsonMap -- contains a deserialized json
type JsonMap struct {
	data interface{}
}

// NewJsonHistoryMap -- Create a HistoryMap from json string content
func NewJsonHistoryMap(data string) (HistoryMap, error) {

	resultMap, err := makeHistoryMapFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &JsonMap{
		data: resultMap,
	}, nil
}

// GetNode -- Get a JSON node from a map structure mapped from a json object or array
func (jm *JsonMap) GetNode(path string) (interface{}, error) {

	var i interface{}
	i = jm.data

	// Special case a root case for just text
	if path == "/" || path == "$" {
		switch t := i.(type) {
		case map[string]interface{}:
			return i, nil
		default:
			return nil, errors.New("Invalid type for / path: " + reflect.TypeOf(t).String())
		}
	}

	// New paths start with $ (only supporting those for now)
	if strings.HasPrefix(path, "$") {
		return jsonpath.Get(path, jm.data)
	} else {
		return getNodeImpl(path, jm.data)
	}
}

func (jm *JsonMap) GetJsonObject() (map[string]interface{}, error) {
	switch t := jm.data.(type) {
	case map[string]interface{}:
		return t, nil
	default:
		return nil, errors.New("Invalid type for Json object")
	}
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
