package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithOutput("join.gif"))
	defer d.Cleanup()

	d.Type(`I=$(gum style I --padding "1 5" --border double --border-foreground 212)`)
	d.Enter()
	d.Type(`LOVE=$(gum style Love --padding "1 4" --border double --border-foreground 57)`)
	d.Enter()
	d.Type(`BUBBLE=$(gum style Bubble --padding "1 8" --border double --border-foreground 255)`)
	d.Enter()
	d.Type(`GUM=$(gum style Gum --padding "1 5" --border double --border-foreground #04B575)`)
	d.Enter()

	time.Sleep(time.Second)
	d.Type("clear", dolly.WithSpeed(50))
	d.Enter()

	d.Type(`TOP=$(gum join  --horizontal "$I" "$LOVE")`)
	d.Enter()
	time.Sleep(time.Second)
	d.Type(`BOTTOM=$(gum join  --horizontal "$BUBBLE" "$GUM")`)
	d.Enter()
	time.Sleep(time.Second)

	d.Type(`gum join --align center --vertical "$TOP" "$BOTTOM"`)
	d.Enter()
	time.Sleep(5 * time.Second)
}
