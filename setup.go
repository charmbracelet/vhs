package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type cleanup func()

func setup() (*rod.Page, cleanup) {
	ttyd := ttydCmd()
	go ttyd.Run()

	browser := rod.New().MustConnect()

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

	go func() {
		counter := 0
		for {
			counter++
			page.MustScreenshot(fmt.Sprintf(capturesPath, counter))
			time.Sleep(time.Second / 60)
		}
	}()

	return page, func() {
		browser.MustClose()
		ttyd.Process.Kill()
	}
}
