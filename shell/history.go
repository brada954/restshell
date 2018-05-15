package shell

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var history = make([]Result, 0)

type Result struct {
	Text       string
	Map        interface{}
	Error      error
	HttpStatus int
	HeaderMap  map[string]string
	AuthMap    interface{}
}

func (r *Result) GetObjectMap() (map[string]interface{}, bool) {
	if m, ok := r.Map.(map[string]interface{}); ok {
		return m, true
	}
	return nil, false
}

func (r *Result) GetArrayMap() ([]interface{}, bool) {
	if m, ok := r.Map.([]interface{}); ok {
		return m, true
	}
	return nil, false
}

var (
	ErrArguments        = errors.New("Invalid arguments")
	ErrInvalidValue     = errors.New("Invalid value type")
	ErrNotFound         = errors.New("Node not found")
	ErrInvalidKey       = errors.New("Node path error")
	ErrUnexpectedType   = errors.New("Node is unexpected type")
	ErrDataType         = errors.New("Invalid history data type")
	ErrNoHistory        = errors.New("History not present")
	ErrArrayOutOfBounds = errors.New("Array index out of bounds")
)

func PushResponse(resp *RestResponse, resperror error) error {
	var result Result

	headers := make(map[string]string, 0)
	for n, values := range resp.GetHeader() {
		headers[n] = values[0]
		if n == "Authorization" {
			authmap, err := decodeJwtClaims(values[0])
			if err == nil {
				if IsCmdDebugEnabled() {
					fmt.Fprintf(ConsoleWriter(), "Pushing AuthMap:\n%v\n", authmap)
				}
				result.AuthMap = authmap
			} else {
				if IsCmdDebugEnabled() {
					fmt.Fprintln(ConsoleWriter(), "Unable to decode JWT")
				}
			}
		}
	}

	result.HeaderMap = headers
	result.Text = resp.Text
	result.Error = resperror
	result.HttpStatus = resp.GetStatus()
	resultMap, err := makeResultMap(resp.Text)
	if err != nil {
		return err
	}
	result.Map = resultMap
	PushResult(result)
	return resperror
}

func PushError(resperror error) error {
	headers := make(map[string]string, 0)

	var result Result
	result.HeaderMap = headers
	result.Text = resperror.Error()
	result.Error = resperror
	result.HttpStatus = -1
	resultMap, err := makeResultMap(result.Text)
	if err != nil {
		return err
	}
	result.Map = resultMap
	PushResult(result)
	return resperror
}

func makeResultMap(data string) (interface{}, error) {
	var f interface{}

	err := json.Unmarshal([]byte(data), &f)
	if err != nil {
		// For unmarshall failures place text at a root node
		return makeRootMap(data)
	}

	if m, ok := f.(map[string]interface{}); ok {
		return m, nil
	}

	if m, ok := f.([]interface{}); ok {
		return m, nil
	}

	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "Unknown data type for history buffer")
	}
	return makeRootMap(data)
}

func makeRootMap(text string) (interface{}, error) {
	m := make(map[string]interface{})
	m["/"] = text
	return m, nil
}

func PushResult(result Result) error {
	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "Pushing result into history")
	}
	history = append(history, result)
	if len(history) > 10 {
		history = history[1:]
	}
	return nil
}

//
// index is a index from the end of the array which was the
// last item appended
//
func PeekResult(index int) (Result, error) {
	if len(history) < 1+index {
		return Result{}, ErrNoHistory
	}
	return history[len(history)-(1+index)], nil
}

func GetValueFromHistory(index int, path string) (string, error) {
	result, err := PeekResult(index)
	if err != nil {
		return "", err
	}

	node, err := GetNode(path, result.Map)
	if err != nil {
		return "", err
	}

	switch t := node.(type) {
	case string:
		return t, nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	default:
		return "", errors.New("Invalid data type found")
	}
}

func GetValueFromAuthHistory(index int, path string) (string, error) {
	result, err := PeekResult(index)
	if err != nil {
		return "", err
	}

	if IsCmdDebugEnabled() {
		fmt.Fprintf(ConsoleWriter(), "AuthMap:\n%v\n", result.AuthMap)
	}

	node, err := GetNode(path, result.AuthMap)
	if err != nil {
		return "", err
	}

	switch t := node.(type) {
	case string:
		return t, nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	default:
		return "", errors.New("Invalid data type found")
	}
}

func GetNode(path string, i interface{}) (interface{}, error) {
	parts := strings.SplitN(path, ".", 2)
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, ErrUnexpectedType
	}

	if len(parts) <= 0 {
		return nil, ErrInvalidKey
	}

	arrIndex := -1
	arrParts := strings.SplitN(parts[0], "[", 2)
	if len(arrParts) > 1 {
		index, err := strconv.Atoi(strings.Trim(arrParts[1], "]"))
		if err != nil {
			return nil, ErrInvalidKey
		}
		arrIndex = index
	}

	if arrIndex >= 0 {
		data := m[arrParts[0]]
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
					return GetNode(parts[1], t[arrIndex])
				}
				return nil, ErrArrayOutOfBounds
			case []interface{}:
				if arrIndex < len(t) {
					return GetNode(parts[1], t[arrIndex])
				}
				return nil, ErrArrayOutOfBounds
			}
			return nil, ErrUnexpectedType
		}
	} else {
		if len(parts) > 1 {
			data := m[parts[0]]
			if data == nil {
				return nil, ErrNotFound
			}
			switch t := m[parts[0]].(type) {
			case map[string]interface{}:
				return GetNode(parts[1], t)
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
				} else {
					return nil, nil
				}
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

func GetNodeAsString(path string, i interface{}) (string, error) {
	n, err := GetNode(path, i)
	if err == nil {
		switch t := n.(type) {
		case string:
			return t, nil
		default:
			err = errors.New("Invalid type for string value")
		}
	}
	return "", err
}

func GetNodeAsTime(path string, i interface{}) (time.Time, error) {
	n, err := GetNode(path, i)
	if err == nil {
		return GetValueAsDate(n)
	}
	return time.Time{}, err
}

func GetNodeAsInt64(path string, i interface{}) (int64, error) {
	n, err := GetNode(path, i)
	if err == nil {
		switch t := n.(type) {
		case float64:
			return int64(t), nil
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		default:
			err = errors.New("Invalid type for date/time value")
		}
	}
	return 0, err
}

func GetNodeAsFloat64(path string, i interface{}) (float64, error) {
	n, err := GetNode(path, i)
	if err == nil {
		switch t := n.(type) {
		case float64:
			return t, nil
		case int64:
			return float64(t), nil
		case int:
			return float64(t), nil
		default:
			err = errors.New("Invalid type for date/time value")
		}
	}
	return 0.0, err
}

func GetValueAsDate(i interface{}) (time.Time, error) {
	switch v := i.(type) {
	case string:
		date, err := time.Parse(time.RFC3339Nano, v)
		savedErr := err
		if err != nil {
			date, err = time.Parse(time.UnixDate, v)
			if err != nil {
				date, err = time.Parse("2006-01-02", v)
				if err != nil {
					date, err = time.Parse(time.StampMilli, v)
				}
			}
		}

		if err != nil {
			return time.Time{}, errors.New(fmt.Sprintf("Value not a date: %s (%s)", v, savedErr.Error()))
		}
		return date, nil
	default:
		return time.Time{}, errors.New(fmt.Sprintf("Invalid type for date check: %v", reflect.TypeOf(i)))
	}
}

func decodeJwtClaims(authToken string) (map[string]interface{}, error) {
	parts := strings.Split(authToken, ".")
	if len(parts) != 3 {
		return nil, errors.New("ERROR: Failed to parse auth token: " + authToken)
	}

	data := decodeString(parts[1])
	var resp = make(map[string]interface{}, 0)
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return nil, errors.New("ERROR DECODING CLAIMS: " + err.Error())
	}
	return resp, nil
}

func decodeString(val string) string {
	s, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(val)
	if err != nil {
		fmt.Fprintln(ErrorWriter(), "Base64 Decoder: ", err.Error())
	}
	return string(s)
}
