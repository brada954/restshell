package shell

import (
	"strings"
)

// SortedStringSlice -- sort a slice of strings
func SortedStringSlice(commands []string) []string {
	var sorted []string = make([]string, 0, len(commands))
	for _, k := range commands {
		var index = 0
		for ; index < len(sorted); index++ {
			if strings.Compare(k, sorted[index]) < 0 {
				break
			}
		}
		sorted = append(sorted, "")
		copy(sorted[index+1:], sorted[index:])
		sorted[index] = k
	}
	return sorted
}

// SortedMapKeys -- Sort a list of keys for a map
func SortedMapKeys(mapData map[string]interface{}) []string {
	var sorted []string = make([]string, 0, len(mapData))
	for k := range mapData {
		var index = 0
		for ; index < len(sorted); index++ {
			if strings.Compare(k, sorted[index]) < 0 {
				break
			}
		}
		sorted = append(sorted, "")
		copy(sorted[index+1:], sorted[index:])
		sorted[index] = k
	}
	return sorted
}

// EscapeStringValue -- Quote a string handling '\"' as well
func EscapeStringValue(input string) string {
	result := make([]rune, 0, 1024)
	for _, c := range input {
		if c == '"' {
			result = append(result, '\\', '"')
		} else if c == '\\' {
			result = append(result, '\\', '\\')
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}
