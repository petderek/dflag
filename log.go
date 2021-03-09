package dflag

import (
	"log"
	"os"
)

type logger interface {
	Printf(format string, v ...interface{})
}

type noopLogger struct{}

func (n *noopLogger) Printf(_ string, _ ...interface{}) {}

func debuglogger() *log.Logger {
	return log.New(os.Stderr, "", 0)
}
