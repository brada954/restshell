package modifiers

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValueModifier -- a function for modifying a value in a chain
type ValueModifier func(i interface{}) (interface{}, error)

// NullModifier -- performs a NOP
func NullModifier(i interface{}) (interface{}, error) {
	return i, nil
}

// NilModifier -- converts a nil to a string
func NilModifier(i interface{}) (interface{}, error) {
	if i == nil {
		return "{nil}", nil
	}
	return i, nil
}

// LengthModifier -- Convert the interface value to a length if valid or return error
func LengthModifier(i interface{}) (interface{}, error) {
	switch v := i.(type) {
	case string:
		return len(v), nil
	case float64:
		// TODO: this is not counting decimal portion
		// (probably because float error can make it really long for some values)
		return getIntLength(int(v)), nil
	case int:
		return getIntLength(v), nil
	case int64:
		return getInt64Length(v), nil
	case map[string]interface{}:
		return len(v), nil
	case []interface{}:
		return len(v), nil
	default:
		return nil, fmt.Errorf("Invalid type (%v) for len()", reflect.TypeOf(i))
	}
}

// ConvertToIntModifier -- A value modifier to make a string or a float64
// an integer (float64's will round down)
// Note: XML floats are strings, need to be converted to float then an int
func ConvertToIntModifier(i interface{}) (interface{}, error) {
	switch v := i.(type) {
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		return i, nil
	case float64:
		return int64(v), nil
	}
	return nil, fmt.Errorf("Invalid type (%v) to make int", reflect.TypeOf(i))
}

// ConvertToFloatModifier -- convert a scaler to a floating value
func ConvertToFloatModifier(i interface{}) (interface{}, error) {
	switch v := i.(type) {
	case string:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	}
	return nil, fmt.Errorf("Invalid type (%v) to make float", reflect.TypeOf(i))
}

func StringToLowerModifier(i interface{}) (interface{}, error) {
	switch v := i.(type) {
	case string:
		return strings.ToLower(v), nil
	}
	return nil, fmt.Errorf("Invalid type make lowercase: %v", reflect.TypeOf(i))
}

func StringToUpperModifier(i interface{}) (interface{}, error) {
	switch v := i.(type) {
	case string:
		return strings.ToUpper(v), nil
	}
	return nil, fmt.Errorf("Invalid type make uppercase: %v", reflect.TypeOf(i))
}

func MakeLengthModifier(vmod ValueModifier) ValueModifier {
	return func(i interface{}) (interface{}, error) {
		v, err := vmod(i)
		if err != nil {
			return v, err
		}
		return LengthModifier(v)
	}
}

func MakeToIntModifier(vmod ValueModifier) ValueModifier {
	return func(i interface{}) (interface{}, error) {
		v, err := vmod(i)
		if err != nil {
			return v, err
		}
		return ConvertToIntModifier(v)
	}
}

func MakeToFloatModifier(vmod ValueModifier) ValueModifier {
	return func(i interface{}) (interface{}, error) {
		v, err := vmod(i)
		if err != nil {
			return v, err
		}
		return ConvertToFloatModifier(v)
	}
}

func MakeStringToLowerModifier(vmod ValueModifier) ValueModifier {
	return func(i interface{}) (interface{}, error) {
		v, err := vmod(i)
		if err != nil {
			return v, err
		}
		return StringToLowerModifier(v)
	}
}

func MakeStringToUpperModifier(vmod ValueModifier) ValueModifier {
	return func(i interface{}) (interface{}, error) {
		v, err := vmod(i)
		if err != nil {
			return v, err
		}
		return StringToUpperModifier(v)
	}
}

func MakeRegExModifier(pattern string, vmod ValueModifier) ValueModifier {
	regexp, regexerr := regexp.Compile(pattern)
	return func(i interface{}) (interface{}, error) {
		newValue, err := vmod(i)
		if err != nil {
			return newValue, err
		}

		if regexerr != nil {
			return newValue, regexerr
		}

		switch v := newValue.(type) {
		case string:
			values := regexp.FindStringSubmatch(v)
			if len(values) == 0 {
				return "", nil
			} else if len(values) > 1 {
				return strings.Join(values[1:], ""), nil
			} else {
				return values[0], nil
			}
		default:
			return nil, errors.New("Invalid type for regexp()")
		}
	}
}

func getLength(i interface{}) int {
	switch v := i.(type) {
	case string:
		return len(v)
	case int:
		return getIntLength(v)
	case int64:
		return getIntLength(int(v))
	}
	return -1
}

func getIntLength(i int) int {
	s := strconv.Itoa(i)
	return len(s)
}

func getInt64Length(i int64) int {
	length := 1
	for ; i >= 10; i = i / 10 {
		length = length + 1
	}
	return length
}
