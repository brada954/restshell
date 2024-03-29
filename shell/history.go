package shell

import (
	"encoding/base64"
	"reflect"
	"strconv"
	"time"

	"errors"
	"fmt"
	"net/http"
	"strings"
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

type HistoryMap interface {
	GetNode(string) (interface{}, error)
}

// HistoryOptions -- Common options for history results
type HistoryOptions struct {
	valueIsResultPath *bool // default path into the history result
	valueIsAuthPath   *bool
	valueIsCookiePath *bool
	valueIsHeaderPath *bool
	valueIsHttpStatus *bool
}

// AddHistoryOptions -- Add options for history payload types
func AddHistoryOptions(set CmdSet, payloadType ...ResultPayloadType) HistoryOptions {
	options := HistoryOptions{}

	if isHistoryOptionsRequested(ResultPath, payloadType) {
		options.valueIsResultPath = set.BoolLong("path", 'p', "Use path/value to reference value in history")
	}

	if isHistoryOptionsRequested(AuthPath, payloadType) {
		options.valueIsAuthPath = set.BoolLong("path-auth", 0, "Use path/value to reference JWT AuthToken value in history")
	}

	if isHistoryOptionsRequested(CookiePath, payloadType) {
		options.valueIsCookiePath = set.BoolLong("path-cookie", 0, "Use path/value to reference Cookie value in history")
	}
	if isHistoryOptionsRequested(HeaderPath, payloadType) {
		options.valueIsHeaderPath = set.BoolLong("path-header", 0, "Use path/value to reference Header value in history")
	}

	return options
}

// IsPathOption -- Is the history path option selected
func (ho HistoryOptions) IsResultPathOption() bool {
	return ho.valueIsResultPath != nil && *ho.valueIsResultPath
}

// IsAuthPath -- Is the Authenticadtion token path option selected to parse JWT
func (ho HistoryOptions) IsAuthPath() bool {
	return ho.valueIsAuthPath != nil && *ho.valueIsAuthPath
}

// IsCookiePath -- Is the cookie path option selected to extract cookie
func (ho HistoryOptions) IsCookiePath() bool {
	return ho.valueIsCookiePath != nil && *ho.valueIsCookiePath
}

// IsHeaderPath -- Is the history path option selected
func (ho HistoryOptions) IsHeaderPath() bool {
	return ho.valueIsHeaderPath != nil && *ho.valueIsHeaderPath
}

// IsHeaderPath -- Is the history path option selected
func (ho HistoryOptions) IsHttpStatusPath() bool {
	return ho.valueIsHttpStatus != nil && *ho.valueIsHttpStatus
}

// IsPathOptionEnabled -- True if any history path option is enabled
func (ho HistoryOptions) IsHistoryPathOptionEnabled() bool {
	if ho.IsResultPathOption() || ho.IsAuthPath() || ho.IsCookiePath() || ho.IsHeaderPath() || ho.IsHttpStatusPath() {
		return true
	}
	return false
}

func (ho *HistoryOptions) SetPathOption(payloadType ResultPayloadType) {

	ho.ClearPathOptions()

	switch payloadType {
	case ResultPath:
		if ho.valueIsResultPath != nil {
			*ho.valueIsResultPath = true
		}
	case AuthPath:
		if ho.valueIsAuthPath != nil {
			*ho.valueIsAuthPath = true
		}
	case CookiePath:
		if ho.valueIsCookiePath != nil {
			*ho.valueIsCookiePath = true
		}
	case HeaderPath:
		if ho.valueIsHeaderPath != nil {
			*ho.valueIsHeaderPath = true
		}
	}
}

func (ho *HistoryOptions) ClearPathOptions() {
	if ho.valueIsResultPath != nil {
		*ho.valueIsResultPath = false
	}
	if ho.valueIsAuthPath != nil {
		*ho.valueIsAuthPath = false
	}
	if ho.valueIsCookiePath != nil {
		*ho.valueIsCookiePath = false
	}
	if ho.valueIsHeaderPath != nil {
		*ho.valueIsHeaderPath = false
	}
}

func (ho HistoryOptions) GetNodeFromHistory(index int, path string) (interface{}, error) {
	result, err := PeekResult(index)
	if err != nil {
		return nil, err
	}

	return ho.GetNode(path, result)

}

func (ho HistoryOptions) GetNode(path string, result Result) (interface{}, error) {
	if ho.IsAuthPath() {
		return result.AuthMap.GetNode(path)
	} else if ho.IsCookiePath() {
		return result.CookieMap.GetNode(path)
	} else if ho.IsHeaderPath() {
		return result.HeaderMap.GetNode(path)
	} else {
		return result.BodyMap.GetNode(path)
	}
}

func isHistoryOptionsRequested(t ResultPayloadType, list []ResultPayloadType) bool {
	var hasAllPaths, hasAltPaths bool

	for _, v := range list {
		if v == AllPaths {
			hasAllPaths = true
		}

		if v == AlternatePaths {
			hasAltPaths = true
		}

		if t == v {
			return true
		}
	}

	if hasAllPaths {
		return true
	}

	if hasAltPaths && t != ResultPath {
		return true
	}
	return false
}

// PushResponse -- Push a RestResponse into the history buffer
func PushResponse(resp *RestResponse, resperror error) error {
	var result Result
	{
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
	}

	PushResult(result)
	return resperror
}

// PushResponse -- Push a RestResponse into the history buffer
func PushText(contentType string, data string, resperror error) error {
	var result Result
	{
		result.Text = data
		result.Error = resperror
		result.HttpStatus = http.StatusOK
		result.HttpStatusString = http.StatusText(http.StatusOK)
		result.addParsedContentToResult(contentType, data)
	}

	PushResult(result)
	return resperror
}

// PushError - push a result that is a single string with the error message
func PushError(resperror error) error {

	var result Result
	{
		emptyMap, _ := NewSimpleHistoryMap(make(map[string]string, 0))
		errorMap, _ := NewTextHistoryMap(resperror.Error())

		result.BodyMap = errorMap
		result.HeaderMap = emptyMap
		result.CookieMap = emptyMap
		result.Text = resperror.Error()
		result.Error = resperror
		result.HttpStatus = -1
		result.HttpStatusString = "No Response Received"
	}

	PushResult(result)
	return resperror
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

// PeekResult - Get a history result using an index. Index from the end of
// the array which was the last item appended
func PeekResult(index int) (Result, error) {
	if len(history) < 1+index {
		return Result{}, ErrNoHistory
	}
	return history[len(history)-(1+index)], nil
}

// ConvertNodeValueTostring -- convert a node result to a string value
func ConvertNodeValueToString(node interface{}) (string, error) {
	switch t := node.(type) {
	case string:
		return t, nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
	case int:
		return strconv.FormatInt(int64(t), 10), nil
	case int32:
		return strconv.FormatInt(int64(t), 10), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	default:
		return "", errors.New("invalid data type found")
	}
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
			return time.Time{}, fmt.Errorf("value not a date: %s (%s)", v, savedErr.Error())
		}
		return date, nil
	default:
		return time.Time{}, fmt.Errorf("invalid type for date check: %v", reflect.TypeOf(i))
	}
}

func decodeJwtClaims(authToken string) (HistoryMap, error) {
	parts := strings.Split(authToken, ".")
	if len(parts) != 3 {
		return nil, errors.New("ERROR: Failed to parse auth token: " + authToken)
	}

	data := decodeString(parts[1])
	h, err := NewJsonHistoryMap(data)
	if err != nil {
		return nil, errors.New("ERROR DECODING CLAIMS: " + err.Error())
	}
	return h, nil
}

func decodeString(val string) string {
	s, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(val)
	if err != nil {
		fmt.Fprintln(ErrorWriter(), "Base64 Decoder: ", err.Error())
	}
	return string(s)
}
