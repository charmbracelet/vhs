package main

import (
	"os"
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	os.Symlink("../../../gum/examples/commit.sh", "commit.sh")
	defer os.Remove("commit.sh")

	d := dolly.New(dolly.WithOutput("commit.gif"), dolly.WithFontSize(34))
	defer d.Cleanup()

	d.Type("./commit.sh")
	d.Enter()

	time.Sleep(time.Second)

	d.Type("↓↓", dolly.WithSpeed(200))
	time.Sleep(time.Second / 2)
	d.Enter()

	time.Sleep(time.Second)

	d.Type("gum", dolly.WithSpeed(125))
	time.Sleep(1 * time.Second)
	d.Enter()
	time.Sleep(1 * time.Second)

	d.Type("Gum is sooo tasty", dolly.WithSpeed(125), dolly.WithVariance(0.5))
	time.Sleep(time.Second / 2)
	d.Enter()
	time.Sleep(1 * time.Second)

	d.Type("I love Bubble Gum.", dolly.WithSpeed(100), dolly.WithVariance(0.5))
	time.Sleep(time.Second / 2)

	d.Enter()
	time.Sleep(time.Second / 8)

	d.Enter()

	time.Sleep(time.Second / 2)

	// Lol, we introduce a typo on purpose and then backspace to fix it.
	d.Type("This commit show \bs just how much I love chewing Bubble Gum!!!", dolly.WithSpeed(125), dolly.WithVariance(0.75))

	time.Sleep(2 * time.Second)

	d.CtrlC()

	time.Sleep(1 * time.Second)

	d.Type("→←", dolly.WithSpeed(200))

	time.Sleep(3 * time.Second)
}
