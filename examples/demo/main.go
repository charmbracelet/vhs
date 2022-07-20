package main

import (
	"time"

	"github.com/charmbracelet/frame"
	"github.com/charmbracelet/frame/ffmpeg"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	page, cleanup := frame.New(frame.WithFontSize(42))
	defer cleanup()
	defer ffmpeg.MakeGIF(ffmpeg.WithOutput("demo.gif")).Run()

	for _, kp := range frame.Type("echo 'Hello, Demo!'") {
		time.Sleep(time.Millisecond * 100)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
	}

	page.Keyboard.Type(input.Enter)
	time.Sleep(time.Second)

	cleanup()
}
