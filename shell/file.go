package shell

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetFileContents(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Validate a file exists in full format or with expected extension added.
// Return the file that was verified that exists or an error.
// Note: extension must be all lower case
func GetValidatedFileName(file string, extension string) (string, error) {
	if len(file) == 0 {
		return "", errors.New("The file was not specified")
	}

	if len(extension) > 0 && !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		if !strings.HasSuffix(strings.ToLower(file), strings.ToLower(extension)) {
			if _, err2 := os.Stat(file + extension); err2 == nil {
				file = file + extension
				err = err2
				if IsCmdDebugEnabled() {
					fmt.Fprintf(ConsoleWriter(), "Appending extension to file name: %s\n", file)
				}
			}
		}

		if os.IsNotExist(err) { // Still not exists
			if IsCmdDebugEnabled() {
				fmt.Fprintf(ConsoleWriter(), "Unable to open file: %s\n", file)
			}
			return "", errors.New("The file does not exist")
		}
	}

	if err != nil {
		return file, errors.New("Error accessing file")
	}
	return file, nil
}
