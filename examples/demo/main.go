package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	// Setup dolly with a larger font size and the output GIF as demo.gif
	d := dolly.New(dolly.WithFontSize(42), dolly.WithOutput("demo.gif"))

	// Defer cleanup which tears down all spawned processes and renders the GIF
	defer d.Cleanup()

	// Type a command
	d.Type("echo 'Hello, 多莉!'", dolly.WithSpeed(100), dolly.WithVariance(0.5))
	d.Enter()

	// Give some buffer time for the GIF
	time.Sleep(time.Second)
}
