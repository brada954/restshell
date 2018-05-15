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
		return err
	}

	aliasStore[key] = command
	return nil
}

func RemoveAlias(key string) error {
	key = strings.TrimSpace(strings.ToUpper(key))
	if _, ok := aliasStore[key]; !ok {
		return errors.New("Key not found")
	}

	delete(aliasStore, key)
	return nil
}

func GetAlias(key string) (string, error) {
	key = strings.TrimSpace(strings.ToUpper(key))
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
