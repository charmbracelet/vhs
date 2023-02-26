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
	logLevelVerbose LogLevel = iota
	// It does not log nothing except the publish shareable URL's
	logLevelQuiet
)

var (
	logger    *log.Logger
	loggerOut io.Writer
)

// initLogger configures logger level
func initLogger(level LogLevel) {
	logger = log.New(os.Stderr, "", 0)

	setLogLevel(level)
}

// setLogLevel modify log level
func setLogLevel(level LogLevel) {
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
