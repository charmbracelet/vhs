package main

import (
	"time"

	"github.com/charmbracelet/dolly"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	d := dolly.New(dolly.WithOutput("customization.gif"), dolly.WithHeight(400), dolly.WithFontSize(30))
	defer d.Cleanup()

	d.Type(`gum input --cursor.foreground "#F4AC45" \`)
	d.Enter()
	d.Type(`--prompt.foreground "#04B575" --prompt "What's up? " \`)
	d.Enter()
	d.Type(`--placeholder "Not much, you?" --value "Not much, you?" \`)
	d.Enter()
	d.Type(`--width 80`)
	d.Enter()
	time.Sleep(time.Second)
	d.WithCtrl(input.KeyA)
	time.Sleep(time.Second)
	d.WithCtrl(input.KeyE)
	time.Sleep(time.Second)
	d.WithCtrl(input.KeyU)
	time.Sleep(2 * time.Second)
}
