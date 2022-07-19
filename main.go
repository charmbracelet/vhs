package main

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

const port = 7681
const width = 800
const height = 400

const capturesPath = "captures/input/input-%02d.png"
const gifPath = "captures/input/input.gif"

func main() {
	ttyd := ttydCmd()
	go ttyd.Run()

	browser := rod.New().MustConnect()

	defer browser.MustClose()
	defer ttyd.Process.Kill()

	// Setup the terminal to match Charm Theme.
	// Includes prompt, theme, font, etc...
	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", port))
	page = page.MustSetViewport(width, height, 1, false)
	page.MustWaitIdle()
	page.MustElement(".xterm").Eval(`this.style.padding = '5em'`)
	page.MustElement(".xterm-viewport").Eval(`this.style.overflow = 'hidden'`)
	page.MustElement("textarea").MustInput("PROMPT='%F{#5a56e0}>%f '").MustType(input.Enter)
	page.MustElement("textarea").MustInput("clear").MustType(input.Enter)
	page.MustWaitIdle()

	// Now do whatever you want to record.
	counter := 0

	page.MustElement("textarea").MustInput("gum input --width 80").MustType(input.Enter)
	page.MustWaitIdle()
	page.MustScreenshot(fmt.Sprintf(capturesPath, counter))

	keypresses := []input.Key{
		input.KeyH,
		input.KeyI,
		input.Space,
		input.KeyG,
		input.KeyU,
		input.KeyM,
		shift(input.Digit1),
	}

	for _, kp := range keypresses {
		page.Keyboard.Type(kp)
		page.MustWaitIdle()
		counter++
		page.MustScreenshot(fmt.Sprintf(capturesPath, counter))
	}
	counter++
	page.MustScreenshot(fmt.Sprintf(capturesPath, counter))

	ffmpegCmd().Run()
}

func shift(k input.Key) input.Key {
	k, _ = k.Shift()
	return k
}
