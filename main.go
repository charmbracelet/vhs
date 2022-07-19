package main

import (
	"github.com/go-rod/rod/lib/input"
	"time"
)

func main() {
	page, cleanup := setup()
	defer cleanup()

	for _, kp := range keys("gum input --width 80") {
		time.Sleep(time.Millisecond * 100)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
	}

	page.Keyboard.Type(input.Enter)
	time.Sleep(time.Second)

	for _, kp := range keys("Hello, Gum!") {
		time.Sleep(time.Millisecond * 100)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
	}

	time.Sleep(time.Second)

	err := ffmpegCmd().Run()
	if err != nil {
		panic(err)
	}
}
