package main

import (
	"time"

	"github.com/charmbracelet/dolly"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	d := dolly.New(
		dolly.WithOutput("repo.gif"),
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
	d.Enter()
	waitBetween()
	for i := 0; i < 20; i++ {
		d.Down()
		waitBetween()
	}
	d.Page.Keyboard.Type(input.Tab)
	waitBetween()
	d.Page.Keyboard.Type(input.Tab)
	waitBetween()

	for i := 0; i < 10; i++ {
		d.Down()
		waitBetween()
	}
	d.Enter()
	waitBetween()
	for i := 0; i < 20; i++ {
		d.Down()
		waitBetween()
	}
	d.Type("h")

	// Goto tags tab
	d.Page.Keyboard.Type(input.Tab)
	waitBetween()
	d.Page.Keyboard.Type(input.Tab)
	waitBetween()
	d.Down()
	waitBetween()
	d.Down()
	waitBetween()
	d.Enter()
	time.Sleep(time.Second * 1)

	d.Down()
	waitBetween()
	d.Down()
	waitBetween()
	d.Type("l")
	waitBetween()
	d.Down()
	waitBetween()
	d.Type("l")
	waitBetween()

	time.Sleep(time.Second * 2)
}
