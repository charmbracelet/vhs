package main

import (
	"os"
	"time"

	"github.com/charmbracelet/dolly"
)

const styleFile = "style.sh"
const styleContent = `#!/bin/bash

gum style --foreground 212 --border-foreground 99 --border double --align center --width 50 --margin "1 2" --padding "2 4" 'Bubble Gum (1¢)' 'So sweet and so fresh!'
`

func main() {
	d := dolly.New(dolly.WithOutput("style.gif"), dolly.WithFontSize(26))
	defer d.Cleanup()

	os.WriteFile(styleFile, []byte(styleContent), 0777)
	defer os.Remove(styleFile)

	d.Type("./style.sh", dolly.WithSpeed(50), dolly.WithVariance(0.25))
	time.Sleep(time.Second / 4)
	d.Enter()

	time.Sleep(5 * time.Second)
}
