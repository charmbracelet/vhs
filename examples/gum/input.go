package main

import (
	"os"
	"time"

	"github.com/charmbracelet/frame/ffmpeg"
	"github.com/charmbracelet/frame/keys"
	"github.com/charmbracelet/frame/setup"
	"github.com/go-rod/rod/lib/input"
)

func gumInput() {
	os.RemoveAll("tmp")

	page, cleanup := setup.Frame(setup.DefaultOptions())
	defer cleanup()
	defer ffmpeg.MakeGIF(ffmpeg.DefaultOptions()).Run()

	for _, kp := range keys.Type("gum input") {
		time.Sleep(100 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(100 * time.Millisecond)
	page.Keyboard.Type(input.Enter)
	time.Sleep(1 * time.Second)

	for _, kp := range keys.Type("Hello, gum!") {
		time.Sleep(100 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(100 * time.Millisecond)
	page.Keyboard.Type(input.Enter)
	time.Sleep(1 * time.Second)
}
