package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithOutput("spin.gif"), dolly.WithFontSize(36), dolly.WithHeight(300))
	defer d.Cleanup()

	d.Type("gum spin --title 'Buying Gum...' -- sleep 5", dolly.WithSpeed(75))
	d.Enter()

	time.Sleep(5 * time.Second)
}
