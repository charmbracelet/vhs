package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	// Setup dolly with a larger font size and the output GIF as demo.gif
	d := dolly.New(dolly.DefaultDollyOptions())

	// Defer cleanup which tears down all spawned processes and renders the GIF
	defer d.Cleanup()

	// Type a command
	d.Type("echo 'Hello, Demo!'", dolly.DefaultTypeOptions)

	// Give some buffer time for the GIF
	time.Sleep(time.Second)
}
