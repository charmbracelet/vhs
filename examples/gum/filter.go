package main

import (
	"os"
	"time"

	"github.com/charmbracelet/frame/ffmpeg"
	"github.com/charmbracelet/frame/keys"
	"github.com/charmbracelet/frame/setup"
	"github.com/go-rod/rod/lib/input"
)

func gumFilter() {
	os.RemoveAll("tmp")

	setupOptions := setup.DefaultOptions()
	setupOptions.FontSize = 42
	page, cleanup := setup.Frame(setupOptions)
	defer cleanup()
	ffmpegOptions := ffmpeg.DefaultOptions()
	ffmpegOptions.Output = "filter.gif"
	defer ffmpeg.MakeGIF(ffmpegOptions).Run()

	for _, kp := range keys.Type("gum filter < flavors.txt > choice.txt") {
		time.Sleep(100 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(100 * time.Millisecond)
	page.Keyboard.Type(input.Enter)
	time.Sleep(1 * time.Second)

	for _, kp := range keys.Type("Cher") {
		time.Sleep(100 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(1 * time.Second)

	for _, kp := range keys.Type("\b\b\b\b") {
		time.Sleep(50 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(1 * time.Second)

	for _, kp := range keys.Type("one") {
		time.Sleep(100 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(1 * time.Second)

	for _, kp := range keys.Type("\b\b\b") {
		time.Sleep(50 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(1 * time.Second)

	page.Keyboard.Type(input.ArrowDown)
	time.Sleep(100 * time.Millisecond)
	page.Keyboard.Type(input.ArrowDown)
	time.Sleep(250 * time.Millisecond)
	page.Keyboard.Type(input.Enter)
	time.Sleep(500 * time.Millisecond)

	for _, kp := range keys.Type("cat choice.txt") {
		time.Sleep(50 * time.Millisecond)
		page.Keyboard.Type(kp)
	}

	time.Sleep(500 * time.Millisecond)
	page.Keyboard.Type(input.Enter)
	time.Sleep(time.Second)
}
