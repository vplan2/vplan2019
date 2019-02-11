// Package logger contains setup und access functions
// as interface for go-logging logger
//   Authors: Ringo Hoffmann
package logger

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

// Setup sets configuration for logger
func Setup(format string, level int) {
	formatter := logging.MustStringFormatter(format)
	logging.SetFormatter(formatter)
	logging.SetLevel(logging.Level(level), "main")
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warning(args ...interface{}) {
	log.Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
