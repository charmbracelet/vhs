package main

import (
	"os"
	"time"

	"github.com/charmbracelet/dolly"
)

const storyFile = "story.txt"

func main() {
	d := dolly.New(dolly.WithOutput("write.gif"), dolly.WithFontSize(40), dolly.WithHeight(600))
	defer d.Cleanup()

	defer os.Remove(storyFile)
	d.Type(`gum write > `+storyFile, dolly.WithSpeed(40))
	d.Enter()
	time.Sleep(time.Second)

	d.Type("Once upon a time")

	time.Sleep(time.Second / 4)
	d.Enter()
	time.Sleep(time.Second / 4)

	d.Type("In a land far, far away...")
	time.Sleep(time.Second / 2)

	d.CtrlC()

	time.Sleep(time.Second)

	d.Type("cat "+storyFile, dolly.WithSpeed(40))
	d.Enter()
	time.Sleep(time.Second)
}
