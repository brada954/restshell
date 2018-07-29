package shell

import (
	"strings"

	"github.com/satori/go.uuid"
)

func init() {
	// Register substitutes
	RegisterSubstitutionHandler("newguid", "newguid", NewGuidSubstitute)
	RegisterSubstitutionHandler("tolower", "tolower", ToLowerSubstitute)
	RegisterSubstitutionHandler("toupper", "toupper", ToUpperSubstitute)
}

// NewGuidSubstitute -- Implementatino of guid substitution
func NewGuidSubstitute(cache interface{}, subname string, fmt string, option string) (value string, data interface{}) {
	var guid uuid.UUID

	if cache == nil {
		var err error
		if guid, err = uuid.NewV4(); err != nil {
			guid = uuid.Nil
		}
	} else {
		guid = cache.(uuid.UUID)
	}

	switch fmt {
	default:
		return guid.String(), guid
	}
}

// ToLowerSubstitute -- returns the ToLower value from options parameter with format
// options to use the option parameter in a variable lookup
func ToLowerSubstitute(cache interface{}, subname string, fmt string, option string) (value string, data interface{}) {
	if cache == nil {
		if fmt == "var" {
			value = GetGlobalString(option)
		} else {
			value = option
		}
	}
	return strings.ToLower(value), value
}

// ToUpperSubstitute -- returns the ToUpper value from options parameter with format
// options to use the option parameter in a variable lookup
func ToUpperSubstitute(cache interface{}, subname string, fmt string, option string) (value string, data interface{}) {
	if cache == nil {
		if fmt == "var" {
			value = GetGlobalString(option)
		} else {
			value = option
		}
	}
	return strings.ToUpper(value), value
}
