package main

import (
	"os"
	"time"

	"github.com/charmbracelet/frame/ffmpeg"
	"github.com/charmbracelet/frame/keys"
	"github.com/charmbracelet/frame/setup"
	"github.com/go-rod/rod/lib/input"
)

func gumWrite() {
	os.RemoveAll("tmp")

	setupOptions := setup.DefaultOptions()
	setupOptions.FontSize = 42
	page, cleanup := setup.Frame(setupOptions)
	defer cleanup()
	ffmpegOptions := ffmpeg.DefaultOptions()
	ffmpegOptions.Output = "write.gif"
	defer ffmpeg.MakeGIF(ffmpegOptions).Run()

	for _, kp := range keys.Type("gum write") {
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
	page.Keyboard.Type(input.Escape)
	time.Sleep(1 * time.Second)
}
