package shell

import (
	"errors"
	"strings"
)

var aliasStore map[string]string = make(map[string]string, 0)

func AddAlias(key string, command string, force bool) error {
	key = strings.TrimSpace(strings.ToUpper(key))
	if len(key) < 1 {
		return errors.New("Invalid alias")
	}

	if _, ok := aliasStore[key]; ok {
		if !force {
			return errors.New("Alias already exists")
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

func RemoveAlias(key string) error {
	key = strings.TrimSpace(strings.ToUpper(key))
	if len(key) == 0 {
		return errors.New("Empty key")
	}

	if _, ok := aliasStore[key]; !ok {
		return errors.New("Key not found")
	}

	delete(aliasStore, key)
	return nil
}

func GetAlias(key string) (string, error) {
	key = strings.TrimSpace(strings.ToUpper(key))
	if len(key) == 0 {
		return "", errors.New("Empty key")
	}

	if value, ok := aliasStore[key]; !ok {
		return "", errors.New("Key not found")
	} else {
		return value, nil
	}
}

func GetAllAliasKeys() []string {
	var list []string = make([]string, 0)
	for key, _ := range aliasStore {
		list = append(list, key)
	}
	return SortedStringSlice(list)
}
