package main

import (
	"io"
	"log"
	"os"
)

var (
	logger    *log.Logger
	loggerOut io.Writer
)

// InitLogger configures logger mode
func InitLogger(mode string) {
	loggerOut = os.Stderr
	logger = log.New(loggerOut, "", 0)

	// Quiet must not log any messages to std output
	if mode == "quiet" {
		loggerOut = io.Discard
		logger.SetOutput(loggerOut)
	}
}
