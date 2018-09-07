package functions

import (
	"strconv"
	"strings"
	"time"

	"github.com/brada954/restshell/shell"
)

func init() {
	shell.RegisterSubstitutionHandler(GetDateDefinition)
	shell.RegisterSubstitutionHandler(SetDateDefinition)
	shell.RegisterSubstitutionHandler(CreateDateDefinition)
}

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
	Name:         "setdate",
	Group:        "date",
	FunctionHelp: "Set a date value equal to the option string",
	Formats: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{Item: "local", Description: "Local time value"},
		shell.SubstitutionItemHelp{Item: "utc", Description: "Utc time value"},
		shell.SubstitutionItemHelp{Item: "unix", Description: "Unix timestamp value"},
	},
	OptionDescription: "",
	Options: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{
			Item:        "2006-01-02T15:04:05",
			Description: "Formatted date string",
		},
	},
	Function: SetDateSubstitute,
}

// CreateDateDefinition --
var CreateDateDefinition = shell.SubstitutionFunction{
	Name:         "createdate",
	Group:        "date",
	FunctionHelp: "Create a date value relative to current time (or unix 0)",
	Formats: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{Item: "local", Description: "Local time value"},
		shell.SubstitutionItemHelp{Item: "utc", Description: "Utc time value"},
		shell.SubstitutionItemHelp{Item: "unix", Description: "Unix timestamp value; value 0"},
	},
	OptionDescription: "",
	Options: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{
			Item:        "{modifier=[-]value;modifer=[-]value}",
			Description: "Add or subtract a given value from the date",
		},
		shell.SubstitutionItemHelp{
			Item:        "S",
			Description: "Add the specified seconds to time",
		},
		shell.SubstitutionItemHelp{
			Item:        "n",
			Description: "Add the specified minutes to time",
		},
		shell.SubstitutionItemHelp{
			Item:        "h",
			Description: "Add the specified hours to time",
		},
		shell.SubstitutionItemHelp{
			Item:        "d",
			Description: "Add the specified days to date",
		},
		shell.SubstitutionItemHelp{
			Item:        "m",
			Description: "Add the specified months to date",
		},
		shell.SubstitutionItemHelp{
			Item:        "y",
			Description: "Add the specified years to the date",
		},
		shell.SubstitutionItemHelp{
			Item:        "t",
			Description: "Truncate the given date/time component to zero (y,m,d,h,n,s)",
		},
	},
	Function: CreateDateSubstitute,
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
		return formatDate(inputTime.Local(), option), inputTime
	}
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

// CreateDateSubstitute -- A function that creates a date from a offset from now
// TODO: started function; need to finish
func CreateDateSubstitute(cache interface{}, subname, format string, option string) (value string, date interface{}) {
	var inputTime = time.Time{}

	if cache == nil {
		switch format {
		case "unix":
			inputTime = time.Unix(0, 0)
		case "utc":
			inputTime = time.Now().UTC()
		case "local":
			inputTime = time.Now()
		default:
			inputTime = time.Now()
		}
	}

	years := 0
	months := 0
	days := 0
	duration := time.Duration(0)

	modifiers := strings.Split(option, ";")
	for _, m := range modifiers {
		parts := strings.SplitN(m, "=", 2)
		if len(parts) != 2 {
			continue
		}
		value, err := strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			value = 0
		}

		ivalue := int(value)
		switch parts[0] {
		case "s":
			duration = duration + (time.Duration(value) * time.Second)
		case "n":
			duration = duration + (time.Duration(value) * time.Minute)
		case "h":
			duration = duration + (time.Duration(value) * time.Hour)
		case "d":
			days = days + ivalue
		case "y":
			years = years + ivalue
		default:
		}
	}

	inputTime = inputTime.AddDate(years, months, days)
	inputTime = inputTime.Add(duration)
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
