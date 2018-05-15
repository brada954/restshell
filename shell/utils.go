package shell

import (
	"strings"
)

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

func SortedMapKeys(mapData map[string]interface{}) []string {
	var sorted []string = make([]string, 0, len(mapData))
	for k, _ := range mapData {
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
