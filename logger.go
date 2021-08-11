package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	LogConsole bool
	LogTime    bool
	Prefix     string

	File *os.File
	s    sync.Mutex
}

func (l *Logger) Log(prefix, level string, data ...interface{}) {
	var text []string
	if l.LogTime {
		text = append(text, time.Now().Format("01/02/2006 15:04:05"))
	}
	if len(l.Prefix) > 0 {
		text = append(text, l.Prefix)
	}
	if len(prefix) > 0 {
		text = append(text, "["+prefix+"]")
	}
	if len(level) > 0 {
		text = append(text, "<"+level+">")
	}

	var t string
	if len(text) > 0 {
		t = strings.Join(text, " ") + ": "
	}
	for _, element := range data {
		t += fmt.Sprintf("%v", element) + " "
	}
	if l.File != nil {
		l.s.Lock()
		_, _ = l.File.WriteString(t + "\n")
		l.s.Unlock()
	}
	if l.LogConsole {
		println(t)
	}
}
