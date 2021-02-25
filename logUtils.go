package utils

import (
	"fmt"
	"time"
)

type Logger struct {
	LogLevel   int
	LoggerName string
}

const LevelError = 0
const LevelInfo = 1
const LevelDebug = 2

func GoLog(args ...interface{}) {
	var string_ string
	for _, element := range args {
		string_ += fmt.Sprintf("%v", element) + " "
	}
	timeString := time.Now().String()[:19]
	fmt.Println(fmt.Sprintf("[%v]", timeString)+":", string_) //aaa
}

func (logger Logger) Debug(data ...interface{}) {
	if logger.LogLevel >= LevelDebug {
		var string_ string
		for _, element := range data {
			string_ += fmt.Sprintf("%v", element) + " "
		}
		GoLog(fmt.Sprintf("[%v] [%v]:", logger.LoggerName, "DEBUG"), string_)
	}
}

func (logger Logger) Error(data ...interface{}) {
	if logger.LogLevel >= LevelError {
		var string_ string
		for _, element := range data {
			string_ += fmt.Sprintf("%v", element) + " "
		}
		GoLog(fmt.Sprintf("[%v] [%v]:", logger.LoggerName, "ERROR"), string_)
	}
}

func (logger Logger) Info(data ...interface{}) {
	if logger.LogLevel >= LevelInfo {
		var string_ string //
		for _, element := range data {
			string_ += fmt.Sprintf("%v", element) + " "
		}
		GoLog(fmt.Sprintf("[%v] [%v]:", logger.LoggerName, "INFO"), string_)
	}
}

//F
