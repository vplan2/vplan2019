// Package logger contains setup und access functions
// as interface for go-logging logger
//   Authors: Ringo Hoffmann
package logger

import (
	"github.com/op/go-logging"
)

const mainLoggerName = "main"

var log = logging.MustGetLogger(mainLoggerName)

// Setup sets configuration for logger
func Setup(format string, level int) {
	formatter := logging.MustStringFormatter(format)
	logging.SetFormatter(formatter)
	logging.SetLevel(logging.Level(level), mainLoggerName)
}

// SetLogLevel sets the log level for the current logger
func SetLogLevel(logLevel int) {
	logging.SetLevel(logging.Level(logLevel), mainLoggerName)
}

func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warning(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
