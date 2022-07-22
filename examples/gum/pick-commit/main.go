package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(dolly.WithOutput("pick-commit.gif"), dolly.WithFontSize(30))

	defer d.Cleanup()

	d.Execute("cd pick-commit/sandbox")
	d.Type(`git log --oneline | gum filter | cut -d' ' -f1`)
	d.Enter()
	time.Sleep(time.Second)

	d.Type("feat: ", dolly.WithSpeed(150))
	time.Sleep(time.Second)
	d.Type("↓↓↓", dolly.WithSpeed(350))
	d.Enter()
	time.Sleep(2 * time.Second)
}
