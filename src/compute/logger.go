package compute

import (
	"log"
	"os"
	"fmt"
)

type Log struct {
	log *log.Logger
}

func Logger() (l *Log){
	return &Log{log: log.New(os.Stderr, "[mule] ", log.Ldate|log.Ltime|log.Lmicroseconds)}
}

func (l *Log) logf(f string, args ...interface{}) {
	l.log.Output(2, fmt.Sprintf(f, args...))
}
