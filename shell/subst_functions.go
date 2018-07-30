package shell

import (
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

func init() {
	// Register substitutes
	RegisterSubstitutionHandler("newguid", "newguid", NewGuidSubstitute)
	RegisterSubstitutionHandler("tolower", "tolower", ToLowerSubstitute)
	RegisterSubstitutionHandler("toupper", "toupper", ToUpperSubstitute)
	RegisterSubstitutionHandler("date", "getdate", GetDateSubstitute)
	RegisterSubstitutionHandler("date", "setdate", SetDateSubstitute)

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

// GetDateSubstitute --
func GetDateSubstitute(cache interface{}, subname string, fmt string, option string) (value string, date interface{}) {
	var inputTime time.Time
	var defaultFmt = "2006-01-02 15:04:05"

	if t, ok := cache.(time.Time); !ok {
		inputTime = time.Now()
	} else {
		inputTime = t
	}

	fmt = strings.ToLower(fmt)
	if len(option) == 0 {
		option = defaultFmt
	}
	switch fmt {
	case "utc":
		return inputTime.UTC().Format(option), inputTime
	case "unix":
		return strconv.FormatInt(inputTime.Unix(), 10), inputTime
	default:
		return inputTime.Format(option), inputTime
	}

	// Scratch work to shape into relative date modifiers ; Need to define parameeters better
	// if len(fmt) > 0 && len(option) > 0 {
	// 	modifier := 0
	// 	if ( v, err := int.Parse(option) ; err == nil {
	// 		modifier = v
	// 	}

	// 	var year, month, day, hour, second int

	// 	switch(fmt) {
	// 	case "Year":
	// 		year = modifier
	// 	case "Month":
	// 		month = modifier
	// 	case "Day":
	// 		day = modifier
	// 	case "Hour":
	// 		hours = modifier
	// 	case "Second":
	// 		seconds = modifier
	// 	}
	// 	if year + month + day > 0 {
	// 		inputTime.AddDate(year, month, day)
	// 	} else if hours + seconds > 0 {
	// 		inputTIme.Add(time.Hour * hours + time.Minute * mins + time.Second * sec)
	// 	}
	// }

	// if cache == nil {
	// 	if len(option) > 0 {
	// 		if t, err := time.Parse("2006-01-02 03:04:05", option); err == nil {
	// 			inputTime = t
	// 		}
	// 	} else {

	// 	}
	// }
}

// SetDateSubstitute -- A function that returns an empty string but sets the date
// value used by the date group functions
func SetDateSubstitute(cache interface{}, subname, fmt string, option string) (value string, date interface{}) {
	var inputTime = time.Time{}

	if len(fmt) == 0 {
		fmt = "2006-01-02T15:04:05"
	}

	if cache == nil {
		if len(option) > 0 {
			if t, err := time.Parse(fmt, option); err == nil {
				inputTime = t
			}
		}
	}
	return "", inputTime
}
