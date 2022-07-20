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
	time.Sleep(2 * time.Second)

	d.Type("57", dolly.WithSpeed(150))
	time.Sleep(time.Second)
	d.Enter()

	time.Sleep(10 * time.Second)

	d.Type("↓↓↑", dolly.WithSpeed(350))
	d.Type("oran", dolly.WithSpeed(150))
	time.Sleep(time.Second)
	d.Enter()

	time.Sleep(3 * time.Second)

	d.Type("↓↑", dolly.WithSpeed(350))
	time.Sleep(time.Second)

	d.Enter()

	time.Sleep(10 * time.Second)
}
