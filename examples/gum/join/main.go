package main

import (
	"os"
	"time"

	"github.com/charmbracelet/dolly"
)

const boxesFile = `border-boxes.sh`
const boxesContent = `#!/bin/bash

export I=$(gum style --padding "1 5" --border double --border-foreground 212 "I")
export LOVE=$(gum style --padding "1 4" --border double --border-foreground 57 "LOVE")
export BUBBLE=$(gum style --padding "1 8" --border double --border-foreground 255 "Bubble")
export GUM=$(gum style --padding "1 5" --border double --border-foreground "#04B575" "Gum")
`

func main() {
	os.WriteFile(boxesFile, []byte(boxesContent), 777)
	defer os.Remove(boxesFile)

	d := dolly.New(dolly.WithOutput("join.gif"), dolly.WithFontSize(24))
	defer d.Cleanup()

	d.Type("source ./"+boxesFile, dolly.WithSpeed(0))
	d.Enter()
	d.Clear()

	d.Type(`TOP=$(gum join  --horizontal "$I" "$LOVE")`)
	d.Enter()
	time.Sleep(time.Second / 2)

	d.Type(`BOTTOM=$(gum join  --horizontal "$BUBBLE" "$GUM")`)
	d.Enter()
	time.Sleep(time.Second / 2)

	d.Type(`gum join --align center --vertical "$TOP" "$BOTTOM"`)
	d.Enter()
	time.Sleep(5 * time.Second)
}
