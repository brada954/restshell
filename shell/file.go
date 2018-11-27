package shell

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// GetFileContentsOfType -- read the file contents using a default extension
// if one was not provided
func GetFileContentsOfType(file string, extension string) (string, error) {
	filename, err := GetValidatedFileName(file, extension)
	if err != nil {
		return "", err
	}

	body, err := GetFileContents(filename)
	if err != nil {
		return "", err
	}
	return body, nil
}

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

// OpenFileForOutput -- open a file
func OpenFileForOutput(name string, truncate bool, append bool, newfile bool) (*os.File, error) {
	var file *os.File
	if _, err := os.Stat(name); err == nil {
		if !(truncate || append) {
			if newfile {
				err = errors.New("Exists")
				for suffix := 1; err.Error() == "Exists"; suffix++ {
					file, err := openNewFileWithSuffix(name, suffix)
					if err != nil && err.Error() != "Exists" {
						return nil, err
					}
					if err == nil {
						return file, nil
					}
				}
			}
			return nil, errors.New("File exists; use --append or --truncate to use the file")
		}
		flags := os.O_APPEND | os.O_WRONLY
		if truncate {
			flags = os.O_WRONLY
		}
		file, err = os.OpenFile(name, flags, 0644)
		if err != nil {
			return nil, errors.New("Open failed: " + err.Error())
		}
		if truncate {
			file.Truncate(0)
		}
	} else {
		file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, errors.New("Open failed: " + err.Error())
		}
	}
	return file, nil
}

func openNewFileWithSuffix(name string, suffix int) (*os.File, error) {
	rootName := name
	ext := ""
	i := strings.LastIndex(name, ".")
	if i >= 0 {
		rootName = name[:i]
		ext = name[i:]
	}

	name = rootName + strconv.Itoa(suffix) + ext
	if _, err := os.Stat(name); err == nil {
		return nil, errors.New("Exists")
	} else {
		var file *os.File
		file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
