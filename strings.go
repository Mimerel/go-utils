package go_utils

import "strings"

/**
Spits the row information in cells
*/
func splitRowValues(row string, seperator string) (parts []string, err error) {
	replacementSeperator := "§"
	row = strings.Replace(row, seperator, replacementSeperator, -1)
	replaceIt := true
	newRow := ""
	for _, char := range row {
		if strings.EqualFold(string(char), "\"") {
			replaceIt = !replaceIt
		}
		if replaceIt && strings.EqualFold(string(char), replacementSeperator) {
			newRow += seperator
		} else {
			newRow += string(char)
		}
	}
	parts = strings.Split(newRow, seperator)
	for k, _ := range parts {
		parts[k] = strings.Replace(parts[k], replacementSeperator, seperator, -1)
		parts[k] = strings.Replace(parts[k], "\"", "", -1)
	}
	return parts, nil
}
