package main

import (
	"time"

	"github.com/go-rod/rod/lib/input"
	"github.com/maaslalani/frame/ffmpeg"
	"github.com/maaslalani/frame/keys"
)

func main() {
	page, cleanup := setup(Options{
		FramePath: "captures/input-%02d.png",
		FrameRate: 60,
		Port:      7681,
		Width:     1200,
		Height:    600,
	})
	defer cleanup()
	defer ffmpeg.MakeGIF(ffmpeg.Options{
		Input:     "captures/input-%02d.png",
		Output:    "captures/input.gif",
		Framerate: 50,
		Width:     1200,
		MaxColors: 256,
	}).Run()

	for _, kp := range keys.Type("gum input --width 80") {
		time.Sleep(time.Millisecond * 100)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
	}

	page.Keyboard.Type(input.Enter)
	time.Sleep(time.Second)

	for _, kp := range keys.Type("Hello, Gum!") {
		time.Sleep(time.Millisecond * 100)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
	}

	time.Sleep(time.Second)

	cleanup()

}
