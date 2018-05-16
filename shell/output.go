package shell

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ShortDisplayFunc -- A function can be used to pretty format the output or condense it
type ShortDisplayFunc func(io.Writer, Result) error

func RestCompletionHandler(response *RestResponse, err error, shortDisplay ShortDisplayFunc) error {

	if err != nil {
		return errors.New("Network Error: " + err.Error())
	}

	if IsCmdDebugEnabled() {
		fmt.Fprintln(ConsoleWriter(), "Displaying response:")
	}

	options := GetDefaultDisplayOptions()
	if IsShort(options) {
		if shortDisplay != nil {
			data, err := PeekResult(0)
			if err != nil {
				return errors.New("Warning: Unable to parse response")
			}
			if data.HttpStatus == http.StatusOK {
				err = shortDisplay(OutputWriter(), data)
			} else {
				options = append(options, Body)
			}
		} else {
			options = append(options, Body)
		}
	}

	response.DumpResponse(OutputWriter(), options...)

	// Return the short error message if not nil
	if err != nil {
		return err
	}

	// Return error for http status errors
	if response.GetStatus() != http.StatusOK {
		return makeStatusError(response.httpResp)
	}
	return nil
}

func makeStatusError(resp *http.Response) error {
	return fmt.Errorf("HTTP Status: %s", resp.Status)
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

type DisplayOption int

const (
	Body DisplayOption = iota
	Headers
	Cookies
	Status
	Short
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
	return result
}

func (resp *RestResponse) DumpCookies(w io.Writer) {
	for _, v := range resp.GetCookies() {
		fmt.Fprintf(w, "Cookie: %s=%s (%v)\n", v.Name, v.Value, v.Expires)
	}
}

func (resp *RestResponse) DumpHeader(w io.Writer) {
	for k, v := range resp.GetHeader() {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}
}

func (resp *RestResponse) DumpResponse(w io.Writer, options ...DisplayOption) {
	if IsStatus(options) && !IsHeaders(options) {
		fmt.Fprintf(w, "HEADER: Status(%s)\n", resp.httpResp.Status)
	}

	if IsHeaders(options) {
		resp.DumpHeader(w)
	}

	if IsCookies(options) {
		resp.DumpCookies(w)
	}

	if IsBody(options) {
		if isStringBinary(resp.Text) {
			fmt.Fprintln(w, "Response contains too many unprintable characters to display")
		} else {
			fmt.Fprintf(w, "Response:\n%s\n", resp.Text)
		}
	}
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

func isStringBinary(text string) bool {
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
