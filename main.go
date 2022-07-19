package main

import (
	"time"

	"github.com/go-rod/rod/lib/input"
)

func main() {
	page, cleanup := setup()
	defer cleanup()

	keypresses := []input.Key{
		input.KeyG,
		input.KeyU,
		input.KeyM,
		input.Space,
		input.KeyI,
		input.KeyN,
		input.KeyP,
		input.KeyU,
		input.KeyT,
		input.Space,
		input.Minus,
		input.Minus,
		input.KeyW,
		input.KeyI,
		input.KeyD,
		input.KeyT,
		input.KeyH,
		input.Space,
		input.Digit8,
		input.Digit0,
		input.Enter,
		input.KeyH,
		input.KeyI,
		input.Space,
		input.KeyG,
		input.KeyU,
		input.KeyM,
		shift(input.Digit1),
		input.Enter,
	}

	for _, kp := range keypresses {
		time.Sleep(100 * time.Millisecond)
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
		if kp == input.Enter {
			time.Sleep(500 * time.Millisecond)
		}
	}

	time.Sleep(time.Second)

	err := ffmpegCmd().Run()
	if err != nil {
		panic(err)
	}
}

func shift(k input.Key) input.Key {
	k, _ = k.Shift()
	return k
}
