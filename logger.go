package main

import (
	"io"
	"log"
	"os"
)

// LogLevel represents custom type for configure
type LogLevel int

const (
	// Default log mode, it log to stdout
	logLevelVerbose LogLevel = 1
	// It does not log nothing except the publish shareable URL's
	logLevelQuiet LogLevel = 2
)

var (
	logger    *log.Logger
	loggerOut io.Writer
)

// InitLogger configures logger level
func InitLogger(level LogLevel) {
	logger = log.New(os.Stderr, "", 0)

	SetLogLevel(level)
}

// SetLogLevel modify log level
func SetLogLevel(level LogLevel) {
	if level == logLevelVerbose {
		setLogLevelVerbose()
	} else if level == logLevelQuiet {
		setLogLevelQuiet()
	}
}

// setLogLevelQuiet configures log level verbose behaviour
func setLogLevelVerbose() {
	loggerOut = os.Stderr
	logger.SetOutput(loggerOut)
}

// setLogLevelQuiet configures log level quiet behaviour
func setLogLevelQuiet() {
	loggerOut = io.Discard
	logger.SetOutput(loggerOut)
}
