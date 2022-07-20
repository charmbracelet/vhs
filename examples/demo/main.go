package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithFontSize(42), dolly.WithOutput("demo.gif"))
	defer d.Cleanup()

	d.Type("echo 'Hello, Demo!'", 50*time.Millisecond)
	d.Enter()

	time.Sleep(time.Second)
}
