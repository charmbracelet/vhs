package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithOutput("pick-tmux-session.gif"), dolly.WithHeight(350), dolly.WithFontSize(30))
	defer d.Cleanup()

	d.Type(`tmux ls -F \#S | gum filter --placeholder "Pick session"`)
	d.Enter()
	time.Sleep(time.Second)
	d.Type("↓↑", dolly.WithSpeed(250))
	time.Sleep(2 * time.Second)
}
