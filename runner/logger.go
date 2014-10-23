package runner

import (
	logPkg "log"
	"os"
)

type logFunc func(string, ...interface{})

var logger = logPkg.New(os.Stderr, "", 0)

func newLogFunc(prefix string) func(string, ...interface{}) {
	return func(format string, v ...interface{}) {
		logger.Printf(format, v...)
	}
}

func fatal(err error) {
	logger.Fatal(err)
}

type appLogWriter struct{}

func (a appLogWriter) Write(p []byte) (n int, err error) {
	appLog(string(p))

	return len(p), nil
}
