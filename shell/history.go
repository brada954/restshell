package shell

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/subchen/go-xmldom"
)

var history = make([]Result, 0)

// Error variables
var (
	ErrArguments         = errors.New("Invalid arguments")
	ErrInvalidValue      = errors.New("Invalid value type")
	ErrNotFound          = errors.New("Node not found")
	ErrInvalidPath       = errors.New("Node path error")
	ErrInvalidKey        = errors.New("Invalid key")
	ErrUnexpectedType    = errors.New("Node is unexpected type")
	ErrDataType          = errors.New("Invalid history data type")
	ErrNoHistory         = errors.New("History not present")
	ErrArrayOutOfBounds  = errors.New("Array index out of bounds")
	ErrInvalidSubCommand = errors.New("Invalid sub-command")
	ErrNotImplemented    = errors.New("Command not implemented")
)

// ResultContentType -- types of result data
type ResultContentType string

// Types of result data
var (
	ResultContentUnknown ResultContentType = "unknown"
	ResultContentXml     ResultContentType = "xml"
	ResultContentJson    ResultContentType = "json"
	ResultContentText    ResultContentType = "text"
	ResultContentHtml    ResultContentType = "html"
	ResultContentCsv     ResultContentType = "csv"
	ResultContentForm    ResultContentType = "form"
	ResultContentBinary  ResultContentType = "binary"
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

	result.Text = resp.Text
	result.Error = resperror
	result.HttpStatus = resp.GetStatus()
	result.HttpStatusString = resp.GetStatusString()
	result.addParsedContentToResult(resp.GetContentType(), resp.Text)

	PushResult(result)
	return resperror
}

//
// PushResponse -- Push a RestResponse into the history buffer
//
func PushText(contentType string, data string, resperror error) error {
	var result Result

	result.Text = data
	result.Error = resperror
	result.HttpStatus = http.StatusOK
	result.HttpStatusString = http.StatusText(http.StatusOK)
	result.addParsedContentToResult(contentType, data)

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
		result.HttpStatusString = "No Response Received"
		result.Map = makeRootMap(result.Text)
	}

	PushResult(result)
	return resperror
}

// getResultTypeFromResponse -- Get the result type
//   xml, json, text, html, css, csv, media, unknown
func getResultTypeFromContentType(contentType string) ResultContentType {
	// Split off parameters
	parts := strings.Split(contentType, ";")
	if len(parts) == 0 {
		return ResultContentUnknown
	}

	contentType = strings.TrimSpace(strings.ToLower(parts[0]))
	if strings.HasPrefix(contentType, "application/xml") ||
		strings.HasSuffix(contentType, "+xml") {
		return ResultContentXml
	} else if strings.HasPrefix(contentType, "application/json") ||
		strings.HasSuffix(contentType, "+json") {
		return ResultContentJson
	} else if strings.HasPrefix(contentType, "text/plain") ||
		strings.HasPrefix(contentType, "text/calendar") ||
		strings.HasPrefix(contentType, "text/css") {
		return ResultContentText
	} else if strings.HasPrefix(contentType, "text/html") {
		return ResultContentHtml
	} else if strings.HasPrefix(contentType, "application/octet-stream") {
		return ResultContentBinary
	} else if strings.Contains(contentType, "text/csv") {
		return ResultContentCsv
	}

	return ResultContentUnknown
}

func makeResultMapFromJson(data string) (interface{}, error) {
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

	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "Unknown/unsupported data type for history buffer")
	}
	return nil, ErrUnexpectedType
}

func makeXMLDOM(data string) (*xmldom.Document, error) {
	wrapper := xmldom.Must(xmldom.ParseXML("<assertwrapper></assertwrapper>"))
	data = strings.TrimSpace(data)
	if doc, err := xmldom.ParseXML(data); err != nil {
		return wrapper, err
	} else {
		wrapper.Root.AppendChild(doc.Root)
		return wrapper, nil
	}
}

func makeRootMap(text string) interface{} {
	m := make(map[string]interface{})
	m["/"] = strings.TrimSpace(text)
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

	if path == "/" {
		return result.Text, nil
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

func GetValueFromCookieHistory(index int, path string) (string, error) {
	result, err := PeekResult(index)
	if err != nil {
		return "", err
	}

	if IsCmdDebugEnabled() {
		fmt.Fprintf(ConsoleWriter(), "Cookie:\n%v\n", result.CookieMap)
	}

	node, err := GetJsonNode(path, convertToJSONMap(result.CookieMap))
	if err != nil {
		return "", err
	}

	switch t := node.(type) {
	case string:
		return t, nil
	default:
		return "", errors.New("Invalid data type found")
	}
}

func GetValueFromHeaderHistory(index int, path string) (string, error) {
	result, err := PeekResult(index)
	if err != nil {
		return "", err
	}

	if IsCmdDebugEnabled() {
		fmt.Fprintf(ConsoleWriter(), "Headers:\n%v\n", result.HeaderMap)
	}

	node, err := GetJsonNode(path, convertToJSONMap(result.HeaderMap))
	if err != nil {
		return "", err
	}

	switch t := node.(type) {
	case string:
		return t, nil
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

func convertToJSONMap(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for n, v := range m {
		result[n] = v
	}
	return result
}
