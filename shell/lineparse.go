package shell

import (
	"regexp"
	"strings"
)

func LineParse(line string) []string {
	var inQuote = false
	var inEscape = false
	var tokens []string = make([]string, 0, 20)
	var token string = ""
	for _, c := range line {
		if inQuote {
			if inEscape {
				inEscape = false
				token = token + string(c)
				continue
			}

			if c == '"' {
				inQuote = false
			} else if c == '\\' {
				inEscape = true
			} else {
				token = token + string(c)
			}
		} else {
			if c == ' ' {
				if len(token) > 0 {
					tokens = append(tokens, token)
					token = ""
				}
			} else if c == '"' {
				inQuote = true
			} else {
				token = token + string(c)
			}
		}
	}
	if len(token) >= 0 {
		tokens = append(tokens, token)
	}
	return tokens
}

func PerformVariableSubstitution(input string) string {
	var replaceStrings []string = make([]string, 0)

	var filter = func(k string, v interface{}) bool {
		if _, ok := v.(string); !ok {
			return false
		}
		return true
	}

	var replaceBuilder = func(kStr string, v interface{}) {
		if rStr, ok := v.(string); ok {
			replaceStrings = append(replaceStrings, "%%"+kStr+"%%", rStr)
		}
	}

	EnumerateGlobals(replaceBuilder, filter)
	r := strings.NewReplacer(replaceStrings...)
	return r.Replace(input)
}

func IsVariableSubstitutionComplete(input string) bool {

	if regx, err := regexp.Compile(`\%\%.*\%\%`); err == nil {
		if regx.MatchString(input) == false {
			return true
		}
	}
	return false // Note: this is returned in error situations as well (requires investigation)
}
