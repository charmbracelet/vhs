package main

import (
	"time"

	"github.com/charmbracelet/frame/ffmpeg"
	"github.com/charmbracelet/frame/keys"
	"github.com/charmbracelet/frame/setup"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	page, cleanup := setup.Frame(setup.Options{FramePath: "tmp/frame-%02d.png", FrameRate: 60, Width: 1200, Height: 600, Port: 7681, FontSize: 42})
	defer cleanup()
	defer ffmpeg.MakeGIF(ffmpeg.Options{Width: 1200, Input: "tmp/frame-%02d.png", Output: "input.gif", Framerate: 50, MaxColors: 256}).Run()

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
