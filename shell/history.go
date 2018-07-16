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

	"github.com/subchen/go-xmldom"
)

var history = make([]Result, 0)

// Result -- a result that can be placed in the history buffer
// and used by assertion handlers
type Result struct {
	Text        string
	Map         interface{}
	XMLDocument *xmldom.Document
	Error       error
	HttpStatus  int
	HeaderMap   map[string]string
	CookieMap   map[string]string
	AuthMap     interface{}
}

// GetObjectMap -- short cut to get an object if the result was a json object
// returns false if it is not a json(-like) object
func (r *Result) GetObjectMap() (map[string]interface{}, bool) {
	if m, ok := r.Map.(map[string]interface{}); ok {
		return m, true
	}
	return nil, false
}

// GetArraytMap -- short cut to get an array if the result was a json array
// returns false if it is not a json(-like) object
func (r *Result) GetArrayMap() ([]interface{}, bool) {
	if m, ok := r.Map.([]interface{}); ok {
		return m, true
	}
	return nil, false
}

func (r *Result) addCookieMap(resp *RestResponse) error {
	cookies := make(map[string]string, 0)
	for _, cookie := range resp.GetCookies() {
		cookies[cookie.Name] = cookie.Value
	}
	r.HeaderMap = cookies
	return nil // TODO: Are there any error conditions
}

func (r *Result) addHeaderMap(resp *RestResponse) error {
	headers := make(map[string]string, 0)
	for n, values := range resp.GetHeader() {
		headers[n] = values[0]

		// Special code that evalutes authorization header as a JWT
		// This may need revisiting to be more acceptable of possibiltiies
		if n == "Authorization" {
			authmap, err := decodeJwtClaims(values[0])
			if err == nil {
				if IsCmdDebugEnabled() {
					fmt.Fprintf(ConsoleWriter(), "Pushing AuthMap:\n%v\n", authmap)
				}
				r.AuthMap = authmap
			} else {
				if IsCmdDebugEnabled() {
					fmt.Fprintln(ConsoleWriter(), "Unable to decode JWT")
				}
			}
		}
	}
	r.HeaderMap = headers
	return nil // TODO: Are there any error conditions
}

// Error variables
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

//
// PushResponse -- Push a RestResponse into the history buffer
//
func PushResponse(resp *RestResponse, resperror error) error {
	var result Result

	if err := result.addCookieMap(resp); err != nil {
		fmt.Fprintf(ConsoleWriter(), "WARNING: parsing cookies returned: %s", err.Error())
	}

	if err := result.addHeaderMap(resp); err != nil {
		fmt.Fprintf(ConsoleWriter(), "WARNING: parsing header returned: %s", err.Error())
	}

	contentType := "application/octet-stream"
	for k, v := range resp.httpResp.Header {
		if strings.ToLower(k) == "content-type" {
			if len(v) > 0 {
				contentType = strings.ToLower(v[0])
			}
		}
	}

	result.Text = resp.Text
	result.Error = resperror
	result.HttpStatus = resp.GetStatus()

	switch contentType {
	case "application/xml":
		{
			doc, err := makeXMLDOM(resp.Text)
			if err != nil {
				result.Map = makeRootMap(resp.Text)
			} else {
				result.XMLDocument = doc
			}

		}
	case "application/json":
		{
			resultMap, err := makeResultMapFromJson(resp.Text)
			if err != nil {
				resultMap = makeRootMap(resp.Text)
			}
			result.Map = resultMap
		}
	case "text/calandar":
		fallthrough
	case "text/css":
		fallthrough
	case "text/html":
		fallthrough
	case "text/plain":
		result.Map = makeRootMap(resp.Text)
	case "text/csv":
		// TODO: future support
		result.Map = makeRootMap(resp.Text)

	default:
		// Make a default text entry
		result.Map = makeRootMap("Unsupported text returned")
	}

	PushResult(result)
	return resperror
}

// PushError - push a result that is a single string with the error message
func PushError(resperror error) error {
	var result Result
	{
		emptyMap := make(map[string]string, 0)

		result.HeaderMap = emptyMap
		result.CookieMap = emptyMap
		result.Text = resperror.Error()
		result.Error = resperror
		result.HttpStatus = -1
		result.Map = makeRootMap(result.Text)
	}

	PushResult(result)
	return resperror
}

func makeResultMapFromJson(data string) (interface{}, error) {
	var f interface{}

	err := json.Unmarshal([]byte(data), &f)
	if err != nil {
		return nil, ErrInvalidValue
	}

	if m, ok := f.(map[string]interface{}); ok {
		return m, nil
	}

	if m, ok := f.([]interface{}); ok {
		return m, nil
	}

	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "Unknown/unsupported data type for history buffer")
	}
	return nil, ErrUnexpectedType
}

func makeXMLDOM(data string) (*xmldom.Document, error) {
	wrapper := xmldom.Must(xmldom.ParseXML("<assertwrapper></assertwrapper>"))
	data = strings.TrimSpace(data)
	if doc, err := xmldom.ParseXML(data); err != nil {
		return wrapper, ErrInvalidValue
	} else {
		wrapper.Root.AppendChild(doc.Root)
		return wrapper, nil
	}
}

func makeRootMap(text string) interface{} {
	m := make(map[string]interface{})
	m["/"] = text
	return m
}

// PushResult -- push a Result structure into the history buffer
func PushResult(result Result) error {
	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "Pushing the result into history")
	}
	history = append(history, result)
	if len(history) > 10 {
		history = history[1:]
	}
	return nil
}

//
// PeekResult - Get a history result using an index. Index from the end of
// the array which was the last item appended
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

	node, err := GetNode(path, result)
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

	node, err := GetJsonNode(path, result.AuthMap)
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

// GetNodeFromXml - given an xpath return the node or nodes returned with
// the inner text
func GetNodeFromXml(path string, doc *xmldom.Document) (result interface{}, rtnerror error) {

	defer func() {
		if r := recover(); r == nil {
			return // Pass-thru existing error code
		} else {
			rtnerror = errors.New("Error with XPATH: " + path)
		}
	}()

	root := doc.Root

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

// GetNode -- From the result look up the value node in the given
// Result. The path is XPath if the result is Xml and "restshell jsonpath"
// if the result is not Xml
func GetNode(path string, result Result) (interface{}, error) {
	if result.XMLDocument != nil {
		return GetNodeFromXml(path, result.XMLDocument)
	} else {
		return GetJsonNode(path, result.Map)
	}
}

// GetJsonNode -- Get a JSON node from a map structure mapped from a json object or array
func GetJsonNode(path string, i interface{}) (interface{}, error) {
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
					return GetJsonNode(parts[1], t[arrIndex])
				}
				return nil, ErrArrayOutOfBounds
			case []interface{}:
				if arrIndex < len(t) {
					return GetJsonNode(parts[1], t[arrIndex])
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
				return GetJsonNode(parts[1], t)
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

func GetJsonNodeAsString(path string, i interface{}) (string, error) {
	n, err := GetJsonNode(path, i)
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

func GetJsonNodeAsTime(path string, i interface{}) (time.Time, error) {
	n, err := GetJsonNode(path, i)
	if err == nil {
		return GetValueAsDate(n)
	}
	return time.Time{}, err
}

func GetJsonNodeAsInt64(path string, i interface{}) (int64, error) {
	n, err := GetJsonNode(path, i)
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

func GetJsonNodeAsFloat64(path string, i interface{}) (float64, error) {
	n, err := GetJsonNode(path, i)
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

// GetValueAsDate -- given a scaler value in an interface convert it
// to a date if it is can be converted
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
