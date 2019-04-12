package go_utils

import "fmt"

func DefaultLogOutput(message string, args ...interface{}) {
	fmt.Printf(message, args)
}

