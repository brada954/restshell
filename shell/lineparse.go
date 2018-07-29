package shell

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
