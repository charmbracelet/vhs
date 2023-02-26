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

	setOutputLevel(output)
}

// setOutputLevel modify log level
func setOutputLevel(output outputMode) {
	if output == outputVerbose {
		setLogLevelVerbose()
	} else if output == outputQuiet {
		setLogLevelQuiet()
	}
}

// setLogLevelQuiet configures log level verbose behaviour
func setLogLevelVerbose() {
	out = os.Stderr
	logger.SetOutput(out)
}

// setLogLevelQuiet configures log level quiet behaviour
func setLogLevelQuiet() {
	out = io.Discard
	logger.SetOutput(out)
}
