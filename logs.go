package go_utils

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	Debug     LogLevel = 2
	Info      LogLevel = 1
	Error     LogLevel = 0
	InfoColor          = "\033[1;34m%s\033[0m"
	//NoticeColor  = "\033[1;36m%s\033[0m"
	//WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor = "\033[1;31m%s\033[0m"
	DebugColor = "\033[0;36m%s\033[0m"
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
		fmt.Printf(InfoColor + time.Now().Format(time.RFC3339) + " - Info : " + computedMessage + " \n")
	}
}

func (l LogParams) Debug(message string, args ...interface{}) {
	if l.level >= 2 {
		computedMessage := fmt.Sprintf(message, args...)
		fmt.Printf(DebugColor + time.Now().Format(time.RFC3339) + " - Debug :  " + computedMessage + " \n")
	}
}

func (l LogParams) Error(message string, args ...interface{}) {
	if l.level >= 0 {
		computedMessage := fmt.Sprintf(message, args...)
		fmt.Printf(ErrorColor + time.Now().Format(time.RFC3339) + " - Error :  " + computedMessage + " \n")
	}
}
