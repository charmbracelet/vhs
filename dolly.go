package dolly

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// Dolly is the object that controls the setup.
type Dolly struct {
	Options *DollyOptions
	Page    *rod.Page
	Start   func()
	Cleanup func()
}

// DollyOptions is the set of options for the setup.
type DollyOptions struct {
	Framerate  float64
	Height     int
	Padding    string
	Width      int
	FontFamily string
	FontSize   int
	LineHeight float64
	Theme      Theme
	GIF        GIFOptions
}

// DefaultDollyOptions returns the default set of options to use for the setup function.
func DefaultDollyOptions() DollyOptions {
	return DollyOptions{
		Framerate:  60,
		Height:     600,
		Width:      1200,
		Padding:    "5em",
		FontFamily: "SF Mono",
		FontSize:   22,
		LineHeight: 1.2,
		Theme:      DefaultTheme,

		GIF: DefaultGIFOptions,
	}
}

// New sets up ttyd and go-rod for recording frames.
func New() Dolly {
	port := randomPort()
	tty := StartTTY(port)
	go tty.Run()

	browser := rod.New().MustConnect()
	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", port))
	page = page.MustWaitLoad()
	page = page.MustWaitIdle()

	opts := DefaultDollyOptions()

	page.Eval(fmt.Sprintf("term.setOption('fontFamily', '%s')", opts.FontFamily))
	page.Eval(fmt.Sprintf("term.setOption('fontSize', '%d')", opts.FontSize))
	page.Eval(fmt.Sprintf("term.setOption('lineHeight', '%f')", opts.LineHeight))
	theme, err := json.Marshal(opts.Theme)
	if err == nil {
		page.Eval(fmt.Sprintf("term.setOption('theme', %s)", theme))
	}
	page.MustElement(".xterm").Eval(fmt.Sprintf(`this.style.padding = '%s'`, opts.Padding))

	return Dolly{
		Options: &opts,
		Page:    page,
		Start: func() {
			page = page.MustSetViewport(opts.Width, opts.Height, 1, false)
			page.MustEval("window.term.fit")

			page.MustElement("body").Eval(`this.style.overflow = 'hidden'`)
			page.MustElement("#terminal-container").Eval(`this.style.overflow = 'hidden'`)
			page.MustElement(".xterm-viewport").Eval(`this.style.overflow = 'hidden'`)

			page.MustElement("textarea").MustInput("PROMPT='%F{#5a56e0}>%f '").MustType(input.Enter)
			page.MustElement("textarea").MustInput("clear").MustType(input.Enter)
			page.MustWaitIdle()

			os.MkdirAll(opts.GIF.InputFolder, os.ModePerm)

			time.Sleep(2500 * time.Millisecond)

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
						os.WriteFile((opts.GIF.InputFolder + "/" + fmt.Sprintf(frameFileFormat, counter)), screenshot, 0644)
					}
					time.Sleep(time.Second / time.Duration(opts.Framerate))
				}
			}()
		},
		Cleanup: func() {
			// Tear down the processes we started.
			browser.MustClose()
			tty.Process.Kill()

			// Make GIF with frames
			err := MakeGIF(opts.GIF).Run()

			// Cleanup frames if we successfully made the GIF.
			if err == nil {
				os.RemoveAll(opts.GIF.InputFolder)
			}
		},
	}
}
