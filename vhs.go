package vhs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	Test          TestOptions
	Video         VideoOptions
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
		Video:         DefaultVideoOptions,
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

			_ = os.MkdirAll(filepath.Dir(opts.Video.Input), os.ModePerm)

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
						_ = os.WriteFile(fmt.Sprintf(opts.Video.Input, counter), screenshot, 0644)
					}
					time.Sleep(time.Second / time.Duration(opts.Framerate))
				}
			}()
		},
		Cleanup: func() {
			// Tear down the processes we started.
			browser.MustClose()
			_ = tty.Process.Kill()

			// Generate the video(s) with the frames.
			var cmds []*exec.Cmd
			cmds = append(cmds, MakeGIF(opts.Video))
			cmds = append(cmds, MakeMP4(opts.Video))
			cmds = append(cmds, MakeWebM(opts.Video))

			for _, cmd := range cmds {
				if cmd == nil {
					continue
				}
				_ = cmd.Run()
			}

			// Cleanup frames if we successfully made the GIF.
			if opts.Video.CleanupFrames {
				os.RemoveAll(opts.Video.Input)
			}
		},
	}
}
