package shell

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ShortDisplayFunc -- A function can be used to pretty format the output or condense it
type ShortDisplayFunc func(io.Writer, Result) error

func RestCompletionHandler(response *RestResponse, resperr error, shortDisplay ShortDisplayFunc) error {

	if resperr != nil {
		PushError(resperr)
		return errors.New("Network Error: " + resperr.Error())
	}

	PushResponse(response, resperr)
	result, err := PeekResult(0)
	if err != nil {
		return errors.New("Error: Unable to get the result")
	}

	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "Displaying response:")
	}

	options := GetDefaultDisplayOptions()
	if IsShort(options) {
		if shortDisplay != nil {
			if result.HttpStatus == http.StatusOK {
				resperr = shortDisplay(OutputWriter(), result)
			} else {
				options = append(options, Body)
			}
		} else {
			options = append(options, Body)
		}
	}

	if filename := GetCmdOutputFileName(); filename != OptionDefaultOutputFile {
		if o, err := OpenFileForOutput(filename, false, false); err != nil {
			return err
		} else {
			defer o.Close()
			result.DumpResult(o, options...)
		}
	} else {
		result.DumpResult(OutputWriter(), options...)
	}

	// Return the short error message if not nil
	if resperr != nil {
		return resperr
	}

	// Return error for http status errors
	if result.HttpStatus != http.StatusOK {
		return fmt.Errorf("HTTP Status: %s", result.HttpStatusString)
	}
	return nil
}

func ColumnizeTokens(tokens []string, columns int, width int) []string {
	var column = 0
	var line = ""
	var result = make([]string, 0)

	if len(tokens) == 0 {
		return result
	}

	for i := 0; i < len(tokens); i++ {
		token := fmt.Sprintf("%-*s", width, tokens[i])
		if column < columns {
			line = line + token
			column = column + 1
		} else {
			result = append(result, line)
			line = token
			column = 1
		}
	}
	result = append(result, line)
	return result
}

// DisplayOption
type DisplayOption int

// DisplayOption values
const (
	Body DisplayOption = iota
	Headers
	Cookies
	Status
	Short
	Pretty
	All
)

func IsBody(l []DisplayOption) bool {
	return isOptionEnabled(l, Body)
}

func IsShort(l []DisplayOption) bool {
	return isOptionEnabled(l, Short)
}

func IsHeaders(l []DisplayOption) bool {
	return isOptionEnabled(l, Headers)
}

func IsCookies(l []DisplayOption) bool {
	return isOptionEnabled(l, Cookies)
}

func IsStatus(l []DisplayOption) bool {
	return isOptionEnabled(l, Headers) || isOptionEnabled(l, Status)
}

func IsPrettyPrint(l []DisplayOption) bool {
	return isOptionEnabled(l, Pretty)
}

func GetDefaultDisplayOptions() []DisplayOption {
	result := make([]DisplayOption, 0)

	if IsCmdSilentEnabled() {
		if IsCmdOutputBodyEnabled() {
			result = append(result, Body)
		}
		if IsCmdOutputShortEnabled() {
			result = append(result, Short)
		}
	} else {
		if IsCmdVerboseEnabled() && !IsCmdOutputShortEnabled() {
			result = append(result, Body)
		} else if IsCmdOutputBodyEnabled() && IsCmdOutputShortEnabled() {
			result = append(result, Short, Body)
		} else if IsCmdOutputBodyEnabled() {
			result = append(result, Body)
		} else {
			// Note: if there is no short handler, it will dump the full body
			result = append(result, Short)
		}
	}

	if (IsCmdVerboseEnabled() && IsCmdDebugEnabled()) || IsCmdOutputHeaderEnabled() {
		result = append(result, Headers)
	}

	if (IsCmdVerboseEnabled() && IsCmdDebugEnabled()) || IsCmdOutputCookieEnabled() {
		result = append(result, Cookies)
	}

	if IsCmdPrettyPrintEnabled() {
		result = append(result, Pretty)
	}
	return result
}

func isOptionEnabled(list []DisplayOption, option DisplayOption) bool {
	for _, a := range list {
		if a == option {
			return true
		} else if a == All {
			return true
		}
	}
	return false
}

func removeOption(list []DisplayOption, option DisplayOption) []DisplayOption {
	result := make([]DisplayOption, 0)
	for _, a := range list {
		if a != option {
			result = append(result, a)
		}
	}
	return result
}

func isRuneBinary(r rune) bool {
	if r == '\r' || r == '\n' || r == '\t' {
		return false
	}

	if r <= 31 {
		return true
	}

	if r >= 128 {
		return true
	}
	return false
}

func IsStringBinary(text string) bool {
	count := 0
	total := 0
	for _, r := range text {
		total++
		if isRuneBinary(r) {
			count++
			if total > 100 && (((count)*100)/total) > 10 {
				return true
			}
		}
	}
	return false
}
