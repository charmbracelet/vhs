package dolly

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/dolly/ffmpeg"
	"github.com/charmbracelet/dolly/ttyd"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// Dolly is the object that controls the setup.
type Dolly struct {
	Page    *rod.Page
	Cleanup func()
}

// Options is the set of options for the setup.
type Options struct {
	Folder    string
	Format    string
	Framerate float64
	Height    int
	Width     int
	TTY       ttyd.Options
}

// DefaultOptions returns the default set of options to use for the setup function.
func DefaultOptions() Options {
	return Options{
		Framerate: 60,
		Folder:    "tmp",
		Format:    "frame-%02d.png",
		Height:    600,
		Width:     1200,
		TTY:       ttyd.DefaultOptions(),
	}
}

// Option is a function that can be used to set options.
type Option func(*Options)

// WithFolder sets the folder where we should save the frames
func WithFolder(folder string) Option {
	return func(o *Options) {
		o.Folder = folder
	}
}

// WithFormat sets the format string for the frames pngs (default: frame-%02d.png)
func WithFormat(format string) Option {
	return func(o *Options) {
		o.Format = format
	}
}

// WithFPS sets the frames per second.
func WithFPS(fps float64) Option {
	return func(o *Options) {
		o.Framerate = fps
	}
}

// WithHeight sets the height of the frame.
func WithHeight(height int) Option {
	return func(o *Options) {
		o.Height = height
	}
}

// WithWidth sets the width of the frame.
func WithWidth(width int) Option {
	return func(o *Options) {
		o.Width = width
	}
}

// WithPort sets the port to use for the setup.
func WithPort(port int) Option {
	return func(o *Options) {
		o.TTY.Port = port
	}
}

// WithFontSize sets the font size for the setup.
func WithFontSize(size int) Option {
	return func(o *Options) {
		o.TTY.FontSize = size
	}
}

// WithFontFamily sets the font family for the setup.
func WithFontFamily(family string) Option {
	return func(o *Options) {
		o.TTY.FontFamily = family
	}
}

// WithLineHeight sets the line height for the setup.
func WithLineHeight(height float64) Option {
	return func(o *Options) {
		o.TTY.LineHeight = height
	}
}

// New sets up ttyd and go-rod for recording frames.
// Returns the set-up rod.Page and a function for cleanup.
func New(opts ...Option) Dolly {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	tty := ttyd.Start(options.TTY)
	go tty.Run()

	os.MkdirAll(options.Folder, os.ModePerm)

	browser := rod.New().MustConnect()

	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", options.TTY.Port))
	page = page.MustSetViewport(options.Width, options.Height, 1, false)
	page.MustWaitIdle()
	page.MustElement(".xterm").Eval(`this.style.padding = '5em'`)
	page.MustElement(".xterm-viewport").Eval(`this.style.overflow = 'hidden'`)
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
					time.Sleep(time.Second / time.Duration(options.Framerate))
					continue
				}
				os.WriteFile((options.Folder + "/" + fmt.Sprintf(options.Format, counter)), screenshot, 0644)
			}
			time.Sleep(time.Second / time.Duration(options.Framerate))
		}
	}()

	return Dolly{
		Page: page,
		Cleanup: func() {
			// Tear down the processes we started.
			browser.MustClose()
			tty.Process.Kill()

			// Make GIF with frames
			err := ffmpeg.MakeGIF(
				ffmpeg.WithFramerate(50),
				ffmpeg.WithInput(options.Folder+"/"+options.Format),
				ffmpeg.WithOutput(options.Folder+"/_out.gif"),
			).Run()

			// Cleanup frames if we successfully made the GIF.
			if err == nil {
				os.RemoveAll(options.Folder)
			}
		},
	}
}
