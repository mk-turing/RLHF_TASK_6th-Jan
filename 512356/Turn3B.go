package main
import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type logMessage func(level, message string)
type weakLogger struct {
	logger logMessage
	refs   int
}

func (w *weakLogger) incRef() {
	w.refs++
}

func (w *weakLogger) decRef() {
	if w.refs > 0 {
		w.refs--
		if w.refs == 0 {
			w.logger = nil
		}
	}
}

func createLogProcessor(logger logMessage, levelLimit int) logMessage {
	totalLogs := 0
	weakLog := &weakLogger{logger: logger}
	weakLog.incRef()

	return func(level, message string) {
		weakLog.incRef()
		defer weakLog.decRef()