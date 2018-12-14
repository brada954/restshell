package modifiers

import (
	"fmt"
	"reflect"
)

// QuoteModifier -- quote a string
func QuoteModifier(i interface{}) (interface{}, error) {
	switch v := i.(type) {
	case string:
		return QuoteString(v), nil
	}
	return nil, fmt.Errorf("Invalid type to quote: %v", reflect.TypeOf(i))
}

// MakeQuoteModifier -- create a modifier for quoting a string
func MakeQuoteModifier(vmod ValueModifier) ValueModifier {
	return func(i interface{}) (interface{}, error) {
		v, err := vmod(i)
		if err != nil {
			return v, err
		}
		return QuoteModifier(v)
	}
}

// QuoteString -- Quote a string handling '\"' as well
func QuoteString(input string) string {
	result := `"`
	for _, c := range input {
		if c == '"' {
			result = result + `\"`
		} else if c == '\\' {
			result = result + string(`\\`)
		} else {
			result = result + string(c)
		}
	}
	return result + `"`
}
