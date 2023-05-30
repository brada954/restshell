package functions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/brada954/restshell/shell"
)

func init() {
	shell.RegisterSubstitutionHandler(GetDateDefinition)
	shell.RegisterSubstitutionHandler(SetDateDefinition)
	shell.RegisterSubstitutionHandler(SetDateOffsetDefinition)
	shell.RegisterSubstitutionHandler(ModifyDateDefinition)
}

// GetDateDefinition -- Display date based on display option provided (default is current time)
var GetDateDefinition = shell.SubstitutionFunction{
	Name:              "getdate",
	Group:             "date",
	FunctionHelp:      "Display formatted date value for key (default is now)",
	FormatDescription: "Format parameter selects interpretation of date/time:",
	Formats: []shell.SubstitutionItemHelp{
		{Item: "local", Description: "Display Local time value"},
		{Item: "utc", Description: "Display UTC time value"},
		{Item: "unix", Description: "Display Unix timestamp"},
	},
	OptionDescription: "Option is the Golang format for date display",
	Options: []shell.SubstitutionItemHelp{
		{Item: "2006-01-02 15:04:05", Description: "Default Golang format for date and time"},
		{Item: "Mon", Description: "Example Golang format for day of week"},
		{Item: "2006", Description: "Example Golang format for year"},
	},
	Function: GetDateSubstitute,
}

// SetDateDefinition -- Sets the cached date value to value provided (default min date)
var SetDateDefinition = shell.SubstitutionFunction{
	Name:         "setdate",
	Group:        "date",
	FunctionHelp: "Set a date value equal to the option string (default to min date)",
	Formats: []shell.SubstitutionItemHelp{
		{Item: "local", Description: "Parse date as Local"},
		{Item: "utc", Description: "Parse date as UTC"},
		{Item: "unix", Description: "Parse date as Unix timestamp"},
	},
	OptionDescription: "The desired date formatted to format string",
	Options: []shell.SubstitutionItemHelp{
		{Item: "2006-01-02T15:04:05", Description: "Default date format string"},
	},
	Function: SetDateSubstitute,
}

// ModifyDateDefinition -- Initialize the date based on an offset of current time
var ModifyDateDefinition = shell.SubstitutionFunction{
	Name:         "moddate",
	Group:        "date",
	FunctionHelp: "Set date value to an offset of current time (default to UTC) (Deprecated for SetDateOffset)",
	Formats: []shell.SubstitutionItemHelp{
		{Item: "local", Description: "Default to Local time value"},
		{Item: "utc", Description: "Default to UTC time value"},
		{Item: "unix", Description: "Default to Unix time value 0"},
	},
	OptionDescription: "Offset options (d=-2;s=+30;t=hns)",
	Options: []shell.SubstitutionItemHelp{
		{Item: "s", Description: "Add the specified seconds to time"},
		{Item: "n", Description: "Add the specified minutes to time"},
		{Item: "h", Description: "Add the specified hours to time"},
		{Item: "d", Description: "Add the specified days to date"},
		{Item: "m", Description: "Add the specified months to date"},
		{Item: "y", Description: "Add the specified years to the date"},
		{Item: "t", Description: "Truncate date/time component(s) to minimum (t=ymdhns)"},
	},
	Function: SetDateOffsetSubstitute,
}

// SetDateOffsetDefinition -- Initialize the date based on an offset of current time
var SetDateOffsetDefinition = shell.SubstitutionFunction{
	Name:         "setdateoffset",
	Group:        "date",
	FunctionHelp: "Set date value to an offset of current time (default to UTC)",
	Formats: []shell.SubstitutionItemHelp{
		{Item: "local", Description: "Default to Local time value"},
		{Item: "utc", Description: "Default to UTC time value"},
		{Item: "unix", Description: "Default to Unix time value 0"},
	},
	OptionDescription: "Offset options (d=-2;s=+30;t=hns)",
	Options: []shell.SubstitutionItemHelp{
		{Item: "s", Description: "Add the specified seconds to time"},
		{Item: "n", Description: "Add the specified minutes to time"},
		{Item: "h", Description: "Add the specified hours to time"},
		{Item: "d", Description: "Add the specified days to date"},
		{Item: "m", Description: "Add the specified months to date"},
		{Item: "y", Description: "Add the specified years to the date"},
		{Item: "t", Description: "Truncate date/time component(s) to minimum (t=ymdhns)"},
	},
	Function: SetDateOffsetSubstitute,
}

// GetDateSubstitute -- Display date based on display option provided (default is current time)
func GetDateSubstitute(cache interface{}, subname string, format string, option string) (value string, date interface{}) {
	var inputTime time.Time
	var defaultFmt = "2006-01-02 15:04:05"

	if cache == nil {
		inputTime = time.Now()
	} else if t, ok := cache.(time.Time); ok {
		inputTime = t
	} else {
		panic("GetDate substitition failure with cached date")
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

// SetDateSubstitute -- Returns empty string but sets the date value used by the date group functions to option string
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
				} else {
					panic(fmt.Sprintf("SetDate substitution failed for invalid option for date format: %s", option))
				}
			}
		case "local":
			if len(option) > 0 {
				if tm, err := time.ParseInLocation(defaultFmt[:minFormatLen], option, time.Local); err == nil {
					inputTime = tm
				} else {
					panic(fmt.Sprintf("SetDate substitution failed for invalid option for date format: %s", option))
				}
			}
		default:
			panic(fmt.Sprintf("SetDate substitution failed for invalid format: %s", format))
		}
	} else {
		return "", cache
	}
	return "", inputTime
}

// SetDateOffsetSubstitute -- A function initializes date with offset from now
func SetDateOffsetSubstitute(cache interface{}, subname, format string, option string) (value string, date interface{}) {
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
		// Cannot change a cached value -- Panic?
		return "", cache
	}

	if tm, err := applyDateModifiers(inputTime, option); err != nil {
		panic(err.Error())
	} else {
		return "", tm
	}
}

func applyDateModifiers(tm time.Time, modifiers string) (time.Time, error) {
	for _, modifier := range strings.Split(modifiers, ";") {
		parts := strings.SplitN(modifier, "=", 2)
		if len(parts) != 2 {
			continue
		}

		component := strings.ToLower(parts[0])
		value, err := strconv.ParseInt(parts[1], 10, 0)
		if component != "t" && err != nil {
			return tm, fmt.Errorf("SetDateOffset substitition failure to parse option value: %s", parts[1])
		}

		switch component {
		case "s":
			tm = tm.Add(time.Duration(value) * time.Second)
		case "n":
			tm = tm.Add(time.Duration(value) * time.Minute)
		case "h":
			tm = tm.Add(time.Duration(value) * time.Hour)
		case "m":
			tm = tm.AddDate(0, int(value), 0)
		case "d":
			tm = tm.AddDate(0, 0, int(value))
		case "y":
			tm = tm.AddDate(int(value), 0, 0)
		case "t":
			// Perform truncating; modify the component to min value
			for _, c := range parts[1] {
				switch c {
				case 'p': // truncate decimal places of second
					tm = tm.Truncate(time.Second)
				case 's':
					seconds := time.Duration(-tm.Second())
					tm = tm.Add(seconds * time.Second)
				case 'n':
					minutes := time.Duration(-tm.Minute())
					tm = tm.Add(minutes * time.Minute)
				case 'h':
					hours := time.Duration(-tm.Hour())
					tm = tm.Add(hours * time.Hour)
				case 'd':
					tm = tm.AddDate(0, 0, 1-tm.Day())
				case 'm':
					tm = tm.AddDate(0, 1-int(tm.Month()), 0)
				case 'y':
					tm = tm.AddDate(1-int(tm.Year()), 0, 0)
				default:
					return tm, fmt.Errorf("SetDateOffset substitution failed for truncate type: %c", c)
				}
			}
		default:
			return tm, fmt.Errorf("SetDateOffset substitution failed for option component: %s", modifier)
		}
	}
	return tm, nil
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
