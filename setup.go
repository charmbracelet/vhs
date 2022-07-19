package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/maaslalani/frame/ttyd"
)

type cleanup func()

// Options is the set of options for the setup.
type Options struct {
	FramePath string
	FrameRate float64
	Height    int
	Width     int
	Port      int
	FontSize  int
}

func setup(opts Options) (*rod.Page, cleanup) {
	tty := ttyd.Start(ttyd.Options{
		Port:       opts.Port,
		FontFamily: "SF Mono",
		FontSize:   opts.FontSize,
		LineHeight: 1.2,
	})
	go tty.Run()

	// Make directory if it doesn't already exist.
	os.MkdirAll(filepath.Dir(fmt.Sprintf(opts.FramePath, 0)), os.ModePerm)

	browser := rod.New().MustConnect()

	// Setup the terminal to match Charm Theme.
	// Includes prompt, theme, font, etc...
	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", opts.Port))
	page = page.MustSetViewport(opts.Width, opts.Height, 1, false)
	page.MustWaitIdle()
	page.MustElement(".xterm").Eval(`this.style.padding = '5em'`)
	page.MustElement(".xterm-viewport").Eval(`this.style.overflow = 'hidden'`)
	page.MustElement("textarea").MustInput("PROMPT='%F{#5a56e0}>%f '").MustType(input.Enter)
	page.MustElement("textarea").MustInput("clear").MustType(input.Enter)
	page.MustWaitIdle()

	// Wait for terminal overlay to disappear.
	// Ideally, we would hide this with JavaScript but it unfortunately does
	// not have a class selector.
	time.Sleep(2 * time.Second)

	go func() {
		counter := 0
		for {
			counter++
			if page != nil {
				screenshot, err := page.Screenshot(false, &proto.PageCaptureScreenshot{})
				if err != nil {
					break
				}
				os.WriteFile(fmt.Sprintf(opts.FramePath, counter), screenshot, 0644)
			}
			time.Sleep(time.Second / time.Duration(opts.FrameRate))
		}
	}()

	return page, func() {
		browser.MustClose()
		tty.Process.Kill()
	}
}
