package log

import (
	"fmt"
)

type simpleLogger struct{}

func newSimpleLogger() *simpleLogger {
	return &simpleLogger{}
}

func (s *simpleLogger) createLogPoint(log logPoint) {
	fmt.Fprintf(log.b, "%s [%-*s] %s:%d:%s() %s\n", log.time.Format("2006-01-02 15:04:05Z07:00"), config.levelPadding, getPrefix(log.level), log.file, log.fileLine, log.funcName, log.msg)

	log.b.WriteString(log.fields.format())
}
