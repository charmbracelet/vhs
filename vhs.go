package vhs

import (
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// VHS is the object that controls the setup.
type VHS struct {
	Options *VHSOptions
	Page    *rod.Page
	Start   func()
	Cleanup func()
}

// VHSOptions is the set of options for the setup.
type VHSOptions struct {
	Framerate     float64
	Height        int
	Padding       string
	Width         int
	FontFamily    string
	FontSize      int
	LetterSpacing float64
	LineHeight    float64
	Theme         Theme
	GIF           GIFOptions
}

// DefaultVHSOptions returns the default set of options to use for the setup function.
func DefaultVHSOptions() VHSOptions {
	return VHSOptions{
		Framerate:     60,
		Height:        600,
		Width:         1200,
		Padding:       "5em",
		FontFamily:    "SF Mono",
		FontSize:      22,
		LetterSpacing: 1.0,
		LineHeight:    1.2,
		Theme:         DefaultTheme,
		GIF:           DefaultGIFOptions,
	}
}

// New sets up ttyd and go-rod for recording frames.
func New() VHS {
	port := randomPort()
	tty := StartTTY(port)
	go tty.Run() //nolint:errcheck

	browser := rod.New().MustConnect()
	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", port))
	opts := DefaultVHSOptions()

	return VHS{
		Options: &opts,
		Page:    page,
		Start: func() {
			page = page.MustSetViewport(opts.Width, opts.Height, 1, false).
				// Let's wait until we can access the window.term variable
				MustWait("() => window.term != undefined")

			page.MustEval("term.fit")
			page.MustWait("() => document.querySelector('.xterm').childElementCount == 3")

			// There is an annoying overlay that displays how large the terminal is, which goes away after
			// two seconds. We could wait those two seconds (i.e. time.Sleep(2 * time.Second)), but to optimize
			// for the user and GIF generating times, we can remove the overlay manually.
			//
			// This is more complicated than it needs to be since the overlay does not have an ID, or class.
			// The correct solution is to use a CSS selector, but that is not supported in the current version.
			//
			// TODO: Add an ID to the overlay in TTYD and then use that here instead.
			// However, for now, we simply check whether the overlay is active by seeing if .xterm has 3 children.
			page.MustEval("() => document.querySelector('.xterm').lastChild.remove()")

			// Apply options to the terminal
			// By this point the setting commands have been executed, so the `opts` struct is up to date.
			page.MustEval(fmt.Sprintf("() => term.setOption('fontSize', '%d')", opts.FontSize))
			page.MustEval(fmt.Sprintf("() => term.setOption('fontFamily', '%s')", opts.FontFamily))
			page.MustEval(fmt.Sprintf("() => term.setOption('letterSpacing', '%f')", opts.LetterSpacing))
			page.MustEval(fmt.Sprintf("() => term.setOption('lineHeight', '%f')", opts.LineHeight))
			page.MustEval(fmt.Sprintf("() => term.setOption('theme', %s)", opts.Theme.String()))
			page.MustElement(".xterm").MustEval(fmt.Sprintf("() => this.style.padding = '%s'", opts.Padding))

			page.MustElement("textarea").MustInput(" fc -p; PROMPT='%F{#5a56e0}>%f '; clear").MustType(input.Enter)
			page.MustElement("body").MustEval("() => this.style.overflow = 'hidden'")
			page.MustElement("#terminal-container").MustEval("() => this.style.overflow = 'hidden'")
			page.MustElement(".xterm-viewport").MustEval("() => this.style.overflow = 'hidden'")

			_ = os.MkdirAll(opts.GIF.InputFolder, os.ModePerm)

			go func() {
				counter := 0
				for {
					counter++
					if page != nil {
						screenshot, err := page.Screenshot(false, &proto.PageCaptureScreenshot{})
						if err != nil {
							time.Sleep(time.Second / time.Duration(opts.Framerate))
							continue
						}
						_ = os.WriteFile((opts.GIF.InputFolder + "/" + fmt.Sprintf(frameFileFormat, counter)), screenshot, 0644)
					}
					time.Sleep(time.Second / time.Duration(opts.Framerate))
				}
			}()
		},
		Cleanup: func() {
			// Tear down the processes we started.
			browser.MustClose()
			_ = tty.Process.Kill()

			// Make GIF with frames
			err := MakeGIF(opts.GIF).Run()

			// Cleanup frames if we successfully made the GIF.
			if err == nil {
				os.RemoveAll(opts.GIF.InputFolder)
			}
		},
	}
}
