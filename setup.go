package dolly

import (
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// Dolly is the object that controls the setup.
type Dolly struct {
	Page    *rod.Page
	Cleanup func()
}

// DollyOptions is the set of options for the setup.
type DollyOptions struct {
	Folder    string
	Format    string
	Output    string
	Framerate float64
	Height    int
	Width     int
	Padding   string

	TTY TTYOptions
	GIF GIFOptions
}

// DefaultDollyOptions returns the default set of options to use for the setup function.
func DefaultDollyOptions() DollyOptions {
	dir := randomDir()
	gifOptions := DefaultGIFOptions
	gifOptions.InputFolder = dir

	return DollyOptions{
		Framerate: 60,
		Folder:    dir,
		Format:    "frame-%02d.png",
		Output:    "_out.gif",
		Height:    600,
		Width:     1200,
		Padding:   "5em",

		TTY: DefaultTTYOptions,
		GIF: gifOptions,
	}
}

// New sets up ttyd and go-rod for recording frames.
func New(opts DollyOptions) Dolly {
	if opts.TTY.Port == 0 {
		opts.TTY.Port = randomPort()
	}

	tty := StartTTY(opts.TTY)
	go tty.Run()

	os.MkdirAll(opts.Folder, os.ModePerm)

	browser := rod.New().MustConnect()

	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", opts.TTY.Port))
	page = page.MustSetViewport(opts.Width, opts.Height, 1, false)
	page = page.MustWaitLoad()
	page = page.MustWaitIdle()
	page.MustElement(".xterm").Eval(fmt.Sprintf(`this.style.padding = '%s'`, opts.Padding))
	page.MustElement("body").Eval(`this.style.overflow = 'hidden'`)
	page.MustElement("#terminal-container").Eval(`this.style.overflow = 'hidden'`)
	page.MustElement(".xterm-viewport").Eval(`this.style.overflow = 'hidden'`)
	// Fit ttyd xterm window to the screen.
	// ttyd stores its xterm instance at `window.term`
	// https://xtermjs.org/docs/api/addons/fit/
	// https://github.com/tsl0922/ttyd/blob/723ae966939527e8db35f27fb69bac0e02860099/html/src/components/terminal/index.tsx#L167-L196
	page.MustEval("window.term.fit")
	page.MustElement("textarea").MustInput("PROMPT='%F{#5a56e0}>%f '").MustType(input.Enter)
	page.MustElement("textarea").MustInput("clear").MustType(input.Enter)
	page.MustWaitIdle()

	time.Sleep(2 * time.Second)

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
				os.WriteFile((opts.Folder + "/" + fmt.Sprintf(opts.Format, counter)), screenshot, 0644)
			}
			time.Sleep(time.Second / time.Duration(opts.Framerate))
		}
	}()

	return Dolly{
		Page: page,
		Cleanup: func() {
			// Tear down the processes we started.
			browser.MustClose()
			tty.Process.Kill()

			// Make GIF with frames
			err := MakeGIF(opts.GIF).Run()

			// Cleanup frames if we successfully made the GIF.
			if err == nil {
				os.RemoveAll(opts.Folder)
			}
		},
	}
}
