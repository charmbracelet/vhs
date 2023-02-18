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
	loggerOut = os.Stderr
	logger = log.New(loggerOut, "", 0)

	// Quiet must not log any messages to std output
	if level == logLevelQuiet {
		loggerOut = io.Discard
		logger.SetOutput(loggerOut)
	}
}
