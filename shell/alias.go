package shell

import (
	"errors"
	"fmt"
	"strings"
)

var aliasStore map[string]string = make(map[string]string, 0)

// AddAlias - Add an aliased command to the library
func AddAlias(key string, command string, force bool) error {
	key = strings.TrimSpace(strings.ToUpper(key))
	if len(key) < 1 {
		return errors.New("invalid alias")
	}

	if _, ok := aliasStore[key]; ok {
		if !force {
			return errors.New("alias already exists")
		}
	}

	if err := validateCmd(command); err != nil {
		parts := strings.SplitN(command, " ", 2)
		cmdName := "Bad Command"
		if len(parts) > 0 {
			cmdName = parts[0]
		}
		return errors.New(cmdName + ": " + err.Error())
	}

	aliasStore[key] = command
	return nil
}

// RemoveAlias - remove an alias from the library
func RemoveAlias(key string) error {
	key = strings.TrimSpace(strings.ToUpper(key))
	if len(key) == 0 {
		return errors.New("empty key")
	}

	if _, ok := aliasStore[key]; !ok {
		return errors.New("key not found")
	}

	delete(aliasStore, key)
	return nil
}

// GetAlias - get the alias command from the library
func GetAlias(key string) (string, error) {
	key = strings.TrimSpace(strings.ToUpper(key))
	if len(key) == 0 {
		return "", errors.New("empty key")
	}

	if value, ok := aliasStore[key]; !ok {
		return "", errors.New("key not found")
	} else {
		return value, nil
	}
}

// GetAllAliasKeys -- gets a list of keys from the library
func GetAllAliasKeys() []string {
	var list []string = make([]string, 0)
	for key := range aliasStore {
		list = append(list, key)
	}
	return SortedStringSlice(list)
}

// ExpandAlias - Expand alias in input string
func ExpandAlias(command string) (string, error) {

	if len(command) == 0 {
		return command, nil
	}

	parts := strings.SplitN(command, " ", 2)
	if len(parts[0]) <= 0 {
		return "", errors.New("unable to determine command to process")
	}

	if strings.Contains(parts[0], "\"") {
		return "", errors.New("command contains illegal characters (\",')")
	}

	if alias, err := GetAlias(parts[0]); err == nil {
		argString := ""
		if len(parts) > 1 {
			argString = parts[1]
		}

		if IsDebugEnabled() {
			fmt.Fprintf(ConsoleWriter(), "Using alias: %s\n", alias)
		}
		if len(argString) > 0 {
			if strings.HasSuffix(alias, "\\") || strings.HasSuffix(alias, "/") {
				return alias + argString, nil
			} else {
				return alias + " " + argString, nil
			}
		} else {
			return alias, nil
		}
	}
	return command, nil
}
