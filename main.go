package main

import (
	"time"

	"github.com/charmbracelet/frame/ffmpeg"
	"github.com/charmbracelet/frame/keys"
	"github.com/charmbracelet/frame/setup"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	setupOptions := setup.DefaultOptions()
	setupOptions.FontSize = 42
	ffmpegOptions := ffmpeg.DefaultOptions()
	ffmpegOptions.Output = "demo.gif"
	page, cleanup := setup.Frame(setupOptions)
	defer cleanup()

	defer ffmpeg.MakeGIF(ffmpegOptions).Run()

	for _, kp := range keys.Type("echo 'Hello, Demo!'") {
		time.Sleep(time.Millisecond * 100)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
	}

	page.Keyboard.Type(input.Enter)
	time.Sleep(time.Second)

	cleanup()
}
