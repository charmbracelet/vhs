package main

import (
	"os"
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	os.Symlink("../../../gum/examples/demo.sh", "demo.sh")
	defer os.Remove("demo.sh")

	d := dolly.New(dolly.WithOutput("demo.gif"), dolly.WithFontSize(32))
	defer d.Cleanup()

	d.Type("./demo.sh")
	d.Enter()

	time.Sleep(time.Second)

	d.Type("Walter")
	time.Sleep(time.Second)
	d.Enter()

	time.Sleep(4 * time.Second)

	d.Type("Nope, sorry!")
	time.Sleep(time.Second / 2)
	d.Enter()
	d.Enter()
	time.Sleep(time.Second / 2)
	d.Type("I don't trust you.")
	d.Enter()

	time.Sleep(1 * time.Second)

	d.CtrlC()

	time.Sleep(2 * time.Second)

	d.Type(" ↓ ↓ ↓", dolly.WithSpeed(350))
	time.Sleep(time.Second)
	d.Enter()

	time.Sleep(7 * time.Second)

	d.Type("li", dolly.WithSpeed(350))
	time.Sleep(1 * time.Second)
	d.Enter()

	time.Sleep(4 * time.Second)

	d.Type("↓↓↑↑", dolly.WithSpeed(350))
	time.Sleep(time.Second)
	d.Enter()

	time.Sleep(12 * time.Second)
}
