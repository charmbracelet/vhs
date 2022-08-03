package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(
		dolly.WithOutput("landing.gif"),
		dolly.WithFontSize(18),
		dolly.WithWidth(1600),
		dolly.WithHeight(900),
		dolly.WithPadding("2em"),
		dolly.WithDebug(),
	)
	defer d.Cleanup()

	d.Type("ssh git.charm.sh", dolly.WithSpeed(60))
	d.Enter()

	time.Sleep(time.Second * 2)

	// Move up/down
	d.Type("↓", dolly.WithSpeed(300), dolly.WithVariance(4), dolly.WithRepeat(4))
	d.Type("↑", dolly.WithSpeed(300), dolly.WithVariance(2), dolly.WithRepeat(3))

	// Toggle help
	d.Type("?", dolly.WithSpeed(1000))

	// Copy command
	d.Type("c", dolly.WithSpeed(2000))

	// Move between tabs
	d.Type("\t", dolly.WithSpeed(500), dolly.WithRepeat(3))
	d.Type("\t", dolly.WithSpeed(800), dolly.WithRepeat(3))
	d.Type("\t", dolly.WithSpeed(3000))
}
