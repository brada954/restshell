package functions

import (
	"strconv"
	"strings"
	"time"

	"github.com/brada954/restshell/shell"
)

// GetDateDefinition --
var GetDateDefinition = shell.SubstitutionFunction{
	Name:              "getdate",
	Group:             "date",
	FunctionHelp:      "Generate a date",
	FormatDescription: "Format Parameter selects type of time:",
	Formats: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{Item: "local", Description: "Local time value"},
		shell.SubstitutionItemHelp{Item: "utc", Description: "Utc time value"},
		shell.SubstitutionItemHelp{Item: "unix", Description: "Unix timestamp value"},
	},
	OptionDescription: "Option controls format of date string",
	Options: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{
			Item:        "{specification}",
			Description: "Golang format string",
		},
		shell.SubstitutionItemHelp{
			Item:        "2006-01-02 15:04:05.000",
			Description: "Example Golang format for date and time",
		},
		shell.SubstitutionItemHelp{
			Item:        "Mon",
			Description: "Example Golang format for day of week",
		},
		shell.SubstitutionItemHelp{
			Item:        "2006",
			Description: "Example Golang format for year",
		},
	},
	Function: GetDateSubstitute,
}

// SetDateDefinition --
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
		return formatDate(inputTime.UTC(), option), inputTime
	case "unix":
		return strconv.FormatInt(inputTime.Unix(), 10), inputTime
	case "local":
		return formatDate(inputTime.Local(), option), inputTime
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

// formatDate -- formats with some special options beyond golang date format string
func formatDate(t time.Time, option string) string {

	switch option {
	case "dayofweek":
		return t.Weekday().String()
	default:
		return t.Format(option)
	}
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
