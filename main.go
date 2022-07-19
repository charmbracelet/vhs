package main

import (
	"time"

	"github.com/go-rod/rod/lib/input"
	"github.com/maaslalani/frame/ffmpeg"
	"github.com/maaslalani/frame/keys"
)

func main() {
	framesPath := "captures/input-%02d.png"
	width := 1200
	page, cleanup := setup(Options{
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
		Output:    "captures/input.gif",
		Framerate: 50,
		Width:     width,
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
