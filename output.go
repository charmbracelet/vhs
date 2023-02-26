package main

import (
	"io"
	"log"
	"os"
)

// outputMode represents how output it behaves
type outputMode int

const (
	// Default output mode, it outputs to stdout
	outputVerbose outputMode = iota
	// Does not output except the publish shareable URL's
	outputQuiet
)

var (
	logger *log.Logger
	out    io.Writer
)

// initOutput configures output mode
func initOutput(output outputMode) {
	logger = log.New(os.Stderr, "", 0)

	setOutputMode(output)
}

// setOutputMode modify log level
func setOutputMode(output outputMode) {
	if output == outputVerbose {
		setOutputModeVerbose()
	} else if output == outputQuiet {
		setOutputModeQuiet()
	}
}

// setLogLevelQuiet configures log level verbose behaviour
func setOutputModeVerbose() {
	out = os.Stderr
	logger.SetOutput(out)
}

// setOutputModeQuiet configures log level quiet behaviour
func setOutputModeQuiet() {
	out = io.Discard
	logger.SetOutput(out)
}
