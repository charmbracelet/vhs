package main

import (
	"time"

	"github.com/charmbracelet/dolly"
)

func main() {
	d := dolly.New(
		dolly.WithOutput("landing.gif"),
		dolly.WithFontSize(18),
		dolly.WithWidth(1600),
		dolly.WithHeight(900),
		dolly.WithPadding("0em"),
		dolly.WithDebug(),
	)
	defer d.Cleanup()

	d.Type("ssh git.charm.sh", dolly.WithSpeed(50))
	d.Enter()

	time.Sleep(time.Second * 2)

	waitBetween := func() {
		time.Sleep(time.Millisecond * 300)
	}

	d.Down()
	waitBetween()
	d.Down()
	waitBetween()
	d.Down()
	waitBetween()
	d.Up()
	waitBetween()
	d.Up()
	waitBetween()
	waitBetween()
	d.Enter()

	time.Sleep(time.Second * 2)
}
