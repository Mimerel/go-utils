package go_utils

import (
	"fmt"
	"strings"
)

/**
Spits the row information in cells
*/
func splitRowValues(row string, seperator string, debug bool) (parts []string, err error) {
	replacementSeperator := " YZYZY "
	skip := 0;
	if debug {
		fmt.Printf("row before: %s\n", row)
	}
	row = strings.Replace(row, seperator, replacementSeperator, -1)
	replaceIt := true
	newRow := ""
	for k, char := range row {
		if skip <= 0 {
			if strings.EqualFold(string(char), "\"") {
				replaceIt = !replaceIt
			}
			findValue := ""
			if k < (len(row) - len(replacementSeperator)) {
				for i := 0; i < len(replacementSeperator); i++ {
					findValue += string(row[k+i])
				}
			}
			if replaceIt &&
				k < (len(row)-len(replacementSeperator)) &&
				strings.EqualFold(findValue, replacementSeperator) {
				newRow += seperator
				skip = len(replacementSeperator) - 1
			} else {
				newRow += string(char)
			}
		} else {
			skip -= 1
		}

	}

	if debug {
		fmt.Printf("row before split: %s\n", newRow)
	}
	parts = strings.Split(newRow, seperator)
	if debug {
		fmt.Printf("Parts before: %v\n", parts)
	}
	for k, _ := range parts {
		parts[k] = strings.Replace(parts[k], replacementSeperator, seperator, -1)
		parts[k] = strings.Replace(parts[k], "\"", "", -1)
	}
	if debug {
		fmt.Printf("Parts after: %v\n", parts)
	}
	return parts, nil
}
