package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithOutput("spin.gif"), dolly.WithFontSize(34))
	defer d.Cleanup()

	d.Type("gum spin -s line --title 'Buying Gum...' sleep 5", dolly.WithSpeed(65))
	d.Enter()
	time.Sleep(3 * time.Second)
}
