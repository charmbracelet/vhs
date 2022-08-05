package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(
		dolly.WithOutput("repo.gif"),
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

	// Choose repository
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(3))
	d.Enter()
	time.Sleep(time.Second * 1)

	// Scroll readme
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(20))
	// Goto commits tab
	d.Type("\t", dolly.WithSpeed(300), dolly.WithRepeat(2))

	// Choose the 10th commit
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(10))
	d.Enter()
	time.Sleep(time.Second * 1)

	// Scroll commit
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(20))
	// Go back to commits log
	d.Type("h", dolly.WithSpeed(300))

	// Goto tags tab
	d.Type("\t", dolly.WithSpeed(300), dolly.WithRepeat(2))

	// Choose the 2nd tag
	d.Type("↓", dolly.WithSpeed(300), dolly.WithRepeat(2))
	d.Enter()
	time.Sleep(time.Second * 1)

	// Select file
	d.Type("↓↓l↓l", dolly.WithSpeed(300))
	time.Sleep(time.Second * 2)
}
