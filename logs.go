package go_utils

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	Debug LogLevel = 2
	Info  LogLevel = 1
	Error LogLevel = 0
)

func DefaultLogOutput(message string, args ...interface{}) {
	fmt.Printf(message, args)
}

type LogParams struct {
	level LogLevel
}

func NewLogger(level LogLevel) LogParams {
	l := LogParams{level: level}
	return l
}

func (l LogParams) Info(message string, args ...interface{}) {
	if l.level >= 1 {
		computedMessage := fmt.Sprintf(message, args...)
		fmt.Printf(time.Now().Format(time.RFC3339)+" - Info : %s \n", computedMessage)
	}
}

func (l LogParams) Debug(message string, args ...interface{}) {
	if l.level >= 0 {
		computedMessage := fmt.Sprintf(message, args...)
		fmt.Printf(time.Now().Format(time.RFC3339)+" - Debug : %s \n", computedMessage)
	}
}

func (l LogParams) Error(message string, args ...interface{}) {
	if l.level >= 2 {
		computedMessage := fmt.Sprintf(message, args...)
		fmt.Printf(time.Now().Format(time.RFC3339)+" - Error : %s \n", computedMessage)
	}
}
