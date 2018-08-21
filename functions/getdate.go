package functions

import (
	"strconv"
	"strings"
	"time"

	"github.com/brada954/restshell/shell"
)

var GetDateDefinition = shell.SubstitutionFunction{
	Name:              "getdate",
	Group:             "date",
	FunctionHelp:      "Generate a date",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetDateSubstitute,
}

var SetDateDefinition = shell.SubstitutionFunction{
	Name:              "setdate",
	Group:             "date",
	FunctionHelp:      "Set a date value",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          SetDateSubstitute,
}

// GetDateSubstitute --
func GetDateSubstitute(cache interface{}, subname string, format string, option string) (value string, date interface{}) {
	var inputTime time.Time
	var defaultFmt = "2006-01-02 15:04:05"

	if t, ok := cache.(time.Time); !ok {
		inputTime = time.Now()
	} else {
		inputTime = t
	}

	format = strings.ToLower(format)
	if len(option) == 0 {
		option = defaultFmt
	}
	switch format {
	case "utc":
		return inputTime.UTC().Format(option), inputTime
	case "unix":
		return strconv.FormatInt(inputTime.Unix(), 10), inputTime
	case "local":
		return inputTime.Local().Format(option), inputTime
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
func SetDateSubstitute(cache interface{}, subname, format string, option string) (value string, date interface{}) {
	var inputTime = time.Time{}
	defaultFmt := "2006-01-02T15:04:05"
	minFormatLen := min(len(defaultFmt), len(option))

	if cache == nil {
		switch format {
		case "unix":
			if len(option) > 0 {
				inputTime = createUnixTimeFromArg(option)
			} else {
				inputTime = time.Unix(0, 0)
			}
		case "utc":
			if len(option) > 0 {
				if tm, err := time.ParseInLocation(defaultFmt[:minFormatLen], option, time.UTC); err == nil {
					inputTime = tm
				}
			}
		case "local":
			if len(option) > 0 {
				if tm, err := time.ParseInLocation(defaultFmt[:minFormatLen], option, time.Local); err == nil {
					inputTime = tm
				}
			}
		default:
		}
	}
	return "", inputTime
}

func createUnixTimeFromArg(input string) time.Time {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return time.Unix(0, 0)
	}
	tm := time.Unix(i, 0)
	return tm
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
