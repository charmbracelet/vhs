package main

import (
	"time"

	"github.com/charmbracelet/dolly"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	d := dolly.New(dolly.WithOutput("choose.gif"), dolly.WithFontSize(40))
	defer d.Cleanup()

	d.Type("gum choose {1..5}", dolly.WithSpeed(50))
	d.Enter()
	time.Sleep(time.Second)

	d.Type("↓↓↓↑↑", dolly.WithSpeed(350))
	time.Sleep(time.Second)

	d.Enter()
	time.Sleep(time.Second)

	d.WithCtrl(input.KeyL)
	time.Sleep(time.Second / 2)

	d.Type("gum choose --limit 2 Banana Cherry Orange", dolly.WithSpeed(50))

	d.Enter()
	time.Sleep(time.Second)

	d.Type("↓ ↓ ↑", dolly.WithSpeed(350))
	time.Sleep(time.Second)

	d.Enter()
	time.Sleep(time.Second)
}
