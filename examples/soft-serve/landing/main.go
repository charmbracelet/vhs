package main

import (
	"time"

	"github.com/charmbracelet/dolly"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	d := dolly.New(
		dolly.WithOutput("landing.gif"),
		dolly.WithFontSize(18),
		dolly.WithWidth(1600),
		dolly.WithHeight(900),
		dolly.WithPadding("2em"),
		dolly.WithDebug(),
	)
	defer d.Cleanup()

	d.Type("ssh git.charm.sh", dolly.WithSpeed(60))
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
	d.Down()
	waitBetween()
	d.Up()
	waitBetween()
	d.Up()
	waitBetween()
	d.Up()
	waitBetween()
	d.Type("?")
	waitBetween()
	d.Type("c")
	time.Sleep(time.Second * 2)
	for i := 0; i < 3; i++ {
		d.Page.Keyboard.Type(input.Tab)
		time.Sleep(time.Millisecond * 500)
	}
	for i := 0; i < 3; i++ {
		d.Page.Keyboard.Type(input.Tab)
		time.Sleep(time.Millisecond * 800)
	}
	d.Page.Keyboard.Type(input.Tab)
	time.Sleep(time.Second * 1)
	time.Sleep(time.Second * 3)
}
