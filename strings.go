package go_utils

import (
	"fmt"
	"strings"
)

/**
Spits the row information in cells
*/
func splitRowValues(row string, seperator string, debug bool) (parts []string, err error) {
	replacementSeperator := "§§§"
	if debug {
		fmt.Printf("row before: %s", row)
	}
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

	if debug {
		fmt.Printf("row before split: %s", row)
	}
	parts = strings.Split(newRow, seperator)
	if debug {
		fmt.Printf("Parts before: %v", parts)
	}
	for k, _ := range parts {
		parts[k] = strings.Replace(parts[k], replacementSeperator, seperator, -1)
		parts[k] = strings.Replace(parts[k], "\"", "", -1)
	}
	if debug {
		fmt.Printf("Parts after: %v", parts)
	}
	return parts, nil
}
