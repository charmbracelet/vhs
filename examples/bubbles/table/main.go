package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithOutput("table.gif"))
	defer d.Cleanup()

	d.Type("./table")
	d.Enter()

	time.Sleep(time.Second)

	d.Type("↓↓↓↓↓↓↓↑↑↑↓↓", dolly.WithSpeed(350))
	d.Enter()

	time.Sleep(time.Second)
}
