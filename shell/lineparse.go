package shell

import (
	"regexp"
)

// LineParse -- Parse a command line into tokens handling
// escape sequences and double quotes
func LineParse(input string) []string {
	inQuote := false
	inEscape := false
	tokens := make([]string, 0, 20)
	token := ""
	returnEmptyToken := false
	for _, c := range input {
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
				if len(token) > 0 || returnEmptyToken {
					tokens = append(tokens, token)
					token = ""
					returnEmptyToken = false
				}
			} else if c == '"' {
				inQuote = true
				returnEmptyToken = true
			} else {
				token = token + string(c)
			}
		}
	}
	if len(token) > 0 || returnEmptyToken {
		tokens = append(tokens, token)
	}
	return tokens
}

// PerformVariableSubstitution -- Perform variable substitution on a string
func PerformVariableSubstitution(input string) string {
	return SubstituteString(input)

	// replaceStrings := make([]string, 0)

	// var filter = func(k string, v interface{}) bool {
	// 	if _, ok := v.(string); !ok {
	// 		return false
	// 	}
	// 	return true
	// }

	// var replaceBuilder = func(kStr string, v interface{}) {
	// 	if rStr, ok := v.(string); ok {
	// 		replaceStrings = append(replaceStrings, "%%"+kStr+"%%", rStr)
	// 	}
	// }

	// EnumerateGlobals(replaceBuilder, filter)
	// r := strings.NewReplacer(replaceStrings...)
	// return r.Replace(input)
}

// IsVariableSubstitutionComplete -- Validate that variable substitution was
// complete (no variable syntax found)
func IsVariableSubstitutionComplete(input string) bool {

	if regx, err := regexp.Compile(`\%\%.*\%\%`); err == nil {
		if regx.MatchString(input) == false {
			return true
		}
	}
	return false // Note: this is returned in error situations as well (requires investigation)
}
