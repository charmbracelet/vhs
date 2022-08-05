package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(
		dolly.WithOutput("demo.gif"),
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

	// Select repository
	d.Type("↓", dolly.WithSpeed(300), dolly.WithVariance(2), dolly.WithRepeat(4))
	d.Type("↑", dolly.WithSpeed(300), dolly.WithVariance(3), dolly.WithRepeat(3))
	d.Enter()
	time.Sleep(time.Second * 1)

	// Scroll readme
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(20))

	// Goto commits tab
	d.Type("\t\t", dolly.WithSpeed(1000))
	// Select the 9th commit
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(9))
	d.Enter()
	time.Sleep(time.Second * 1)

	// Scroll commit
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(30))
	// Go back to commits log
	d.Type("h", dolly.WithSpeed(300))
	// Goto branches tab
	d.Type("\t", dolly.WithSpeed(300))

	// Select the 2nd branch
	d.Type("↓↓", dolly.WithSpeed(300))
	d.Enter()
	time.Sleep(time.Second * 1)

	// Select file
	d.Type("↓↓", dolly.WithSpeed(300))
	d.Type("ll", dolly.WithSpeed(300))
	// Scroll file
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(10))
}
