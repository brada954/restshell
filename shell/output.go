package shell

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func RestCompletionHandler(response *RestResponse, err error, shortDisplay func(Result) error) error {
	if err != nil {
		return errors.New("Network Error: " + err.Error())
	}

	options := GetDefaultDisplayOptions()
	if !IsCmdSilentEnabled() {
		if !IsCmdVerboseEnabled() && shortDisplay != nil {
			data, err := PeekResult(0)
			if err != nil {
				return errors.New("Warning: Unable to parse response")
			}
			if data.HttpStatus == http.StatusOK {
				return shortDisplay(data)
			} else {
				response.DumpResponse(OutputWriter(), options...)
			}
		} else {
			response.DumpResponse(OutputWriter(), options...)
		}
	}

	if response.GetStatus() != http.StatusOK {
		return statusError(response.GetStatus())
	}
	return nil
}

func statusError(status int) error {
	return errors.New(fmt.Sprintf("Http Failure: %d", status))
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
	All
)

func IsBody(l []DisplayOption) bool {
	return isOptionEnabled(l, Body)
}

func IsHeaders(l []DisplayOption) bool {
	return isOptionEnabled(l, Headers)
}

func IsCookies(l []DisplayOption) bool {
	return isOptionEnabled(l, Cookies)
}

func IsStatus(l []DisplayOption) bool {
	return isOptionEnabled(l, Headers) || isOptionEnabled(l, Body)
}

func GetDefaultDisplayOptions() []DisplayOption {
	result := make([]DisplayOption, 0)

	if !IsCmdSilentEnabled() {
		result = append(result, Body)
	}

	if IsCmdVerboseEnabled() && IsCmdDebugEnabled() {
		result = append(result, Cookies, Headers)
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
	if IsStatus(options) {
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
