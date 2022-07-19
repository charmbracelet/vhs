package main

import (
	"time"

	"github.com/charmbracelet/frame/ffmpeg"
	"github.com/charmbracelet/frame/keys"
	"github.com/charmbracelet/frame/setup"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	framesPath := "tmp/frame-%02d.png"
	width := 1200
	page, cleanup := setup.Frame(setup.Options{
		FramePath: framesPath,
		FrameRate: 60,
		Port:      7681,
		Width:     width,
		Height:    600,
		FontSize:  42,
	})
	defer cleanup()
	defer ffmpeg.MakeGIF(ffmpeg.Options{
		Input:     framesPath,
		Output:    "demo.gif",
		Framerate: 50,
		Width:     width,
		MaxColors: 256,
	}).Run()

	for _, kp := range keys.Type("echo 'Hello, Demo!'") {
		time.Sleep(time.Millisecond * 100)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
	}

	page.Keyboard.Type(input.Enter)
	time.Sleep(time.Second)

	cleanup()
}
