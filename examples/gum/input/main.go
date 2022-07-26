package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithOutput("input.gif"), dolly.WithFontSize(38), dolly.WithHeight(400))
	defer d.Cleanup()

	d.Type(`gum input --placeholder "What's up?"`, dolly.WithSpeed(40))
	d.Enter()
	time.Sleep(time.Second)

	d.Type("Not much, you?")
	d.Enter()
	time.Sleep(time.Second)
}
