package compute

import (
	"fmt"
	"io"
	"log"
)

type Log struct {
	log *log.Logger
}

func Logger(out io.Writer) (l *Log) {
	return &Log{log: log.New(out, "[mule] ", log.Ldate|log.Ltime|log.Lmicroseconds)}
}

func (l *Log) logf(f string, args ...interface{}) {
	l.log.Output(2, fmt.Sprintf(f, args...))
}
