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
	shell.RegisterSubstitutionHandler(ModifyDateDefinition)
}

// GetDateDefinition --
var GetDateDefinition = shell.SubstitutionFunction{
	Name:              "getdate",
	Group:             "date",
	FunctionHelp:      "Return formatted data value (default is now)",
	FormatDescription: "Format Parameter selects type of time:",
	Formats: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{Item: "local", Description: "Display Local time value"},
		shell.SubstitutionItemHelp{Item: "utc", Description: "Display UTC time value"},
		shell.SubstitutionItemHelp{Item: "unix", Description: "Display Unix timestamp"},
	},
	OptionDescription: "Option is the Golang format for a date string",
	Options: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{
			Item:        "2006-01-02 15:04:05",
			Description: "Default Golang format for date and time",
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
	FunctionHelp: "Set a date value equal to the option string (default to min date)",
	Formats: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{Item: "local", Description: "Parse date as Local"},
		shell.SubstitutionItemHelp{Item: "utc", Description: "Parse date as UTC"},
		shell.SubstitutionItemHelp{Item: "unix", Description: "Parse date as Unix timestamp"},
	},
	OptionDescription: "The desired date formatted to format string",
	Options: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{
			Item:        "2006-01-02T15:04:05",
			Description: "Default date format string",
		},
	},
	Function: SetDateSubstitute,
}

// ModifyDateDefinition --
var ModifyDateDefinition = shell.SubstitutionFunction{
	Name:         "moddate",
	Group:        "date",
	FunctionHelp: "Modify current time by component (default to UTC)",
	Formats: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{Item: "local", Description: "Default to Local time value"},
		shell.SubstitutionItemHelp{Item: "utc", Description: "Default to UTC time value"},
		shell.SubstitutionItemHelp{Item: "unix", Description: "Default to Unix time value 0"},
	},
	OptionDescription: "Modifications options (d=-2;s=+30;t=hns)",
	Options: []shell.SubstitutionItemHelp{
		shell.SubstitutionItemHelp{
			Item:        "s",
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
			Description: "Truncate date/time component(s) to minimum (t=ymdhns)",
		},
	},
	Function: ModifyDateSubstitute,
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
	} else {
		return "", cache
	}
	return "", inputTime
}

// ModifyDateSubstitute -- A function that modifies a date from from current value
// Defaults to Now()
func ModifyDateSubstitute(cache interface{}, subname, format string, option string) (value string, date interface{}) {
	var inputTime = time.Time{}

	if cache == nil {
		switch format {
		case "unix":
			inputTime = time.Now().UTC()
			inputTime = inputTime.Truncate(time.Second)
		case "utc":
			inputTime = time.Now().UTC()
		case "local":
			inputTime = time.Now()
		default:
			inputTime = time.Now()
		}
	} else {
		return "", cache
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

		component := strings.ToLower(parts[0])
		value, err := strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			value = 0
		}
		ivalue := int(value)

		switch component {
		case "s":
			duration = duration + (time.Duration(value) * time.Second)
		case "n":
			duration = duration + (time.Duration(value) * time.Minute)
		case "h":
			duration = duration + (time.Duration(value) * time.Hour)
		case "m":
			months = months + ivalue
		case "d":
			days = days + ivalue
		case "y":
			years = years + ivalue
		case "t":
			// Apply current modifications prior to truncating and
			inputTime = inputTime.AddDate(years, months, days)
			inputTime = inputTime.Add(duration)

			// reset the modifiers to zero for continued manipulation
			years = 0
			months = 0
			days = 0
			duration = time.Duration(0)

			// Perform truncating; modify the component to min value
			for _, c := range parts[1] {
				switch c {
				case 'p': // truncate decimal places of second
					inputTime = inputTime.Truncate(time.Second)
				case 's':
					seconds := time.Duration(-inputTime.Second())
					inputTime = inputTime.Add(seconds * time.Second)
				case 'n':
					minutes := time.Duration(-inputTime.Minute())
					inputTime = inputTime.Add(minutes * time.Minute)
				case 'h':
					hours := time.Duration(-inputTime.Hour())
					inputTime = inputTime.Add(hours * time.Hour)
				case 'd':
					inputTime = inputTime.AddDate(0, 0, 1-inputTime.Day())
				case 'm':
					inputTime = inputTime.AddDate(0, 1-int(inputTime.Month()), 0)
				case 'y':
					inputTime = inputTime.AddDate(1-int(inputTime.Year()), 0, 0)
				default:
				}
			}

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
