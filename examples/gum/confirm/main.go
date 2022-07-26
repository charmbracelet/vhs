package main

import (
	"time"

	"github.com/charmbracelet/dolly"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	d := dolly.New(dolly.WithOutput("confirm.gif"), dolly.WithFontSize(32), dolly.WithHeight(400))
	defer d.Cleanup()

	d.Type(`gum confirm && echo 'Me too!' || echo 'Me neither.'`)
	d.Enter()
	time.Sleep(1 * time.Second)

	d.Type("→←", dolly.WithSpeed(300))
	time.Sleep(1 * time.Second)
	d.Enter()
	time.Sleep(2 * time.Second)

	d.Clear()

	time.Sleep(1 * time.Second)

	d.WithCtrl(input.KeyP)
	time.Sleep(1 * time.Second)
	d.Enter()
	time.Sleep(1 * time.Second)

	d.Type("→", dolly.WithSpeed(300))
	time.Sleep(1 * time.Second)
	d.Enter()

	time.Sleep(3 * time.Second)
}
