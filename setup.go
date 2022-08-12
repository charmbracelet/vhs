package dolly

import (
	"fmt"
	"net"
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

// SetupOptions is the set of options for the setup.
type SetupOptions struct {
	Folder    string
	Format    string
	Output    string
	Framerate float64
	Height    int
	Width     int
	Padding   string
	TTY       ttyd.Options
}

// DefaultOptions returns the default set of options to use for the setup function.
func DefaultOptions() SetupOptions {
	tmp, _ := os.MkdirTemp(os.TempDir(), "dolly")
	return SetupOptions{
		Framerate: 60,
		Folder:    tmp,
		Format:    "frame-%02d.png",
		Output:    "_out.gif",
		Height:    600,
		Width:     1200,
		Padding:   "5em",
		TTY:       ttyd.DefaultOptions(),
	}
}

// SetupOption is a function that can be used to set options.
type SetupOption func(*SetupOptions)

// WithFolder sets the folder where we should save the frames
func WithFolder(folder string) SetupOption {
	return func(o *SetupOptions) {
		o.Folder = folder
	}
}

// WithFormat sets the format string for the frames pngs (default: frame-%02d.png)
func WithFormat(format string) SetupOption {
	return func(o *SetupOptions) {
		o.Format = format
	}
}

// WithFPS sets the frames per second.
func WithFPS(fps float64) SetupOption {
	return func(o *SetupOptions) {
		o.Framerate = fps
	}
}

// WithHeight sets the height of the frame.
func WithHeight(height int) SetupOption {
	return func(o *SetupOptions) {
		o.Height = height
	}
}

// WithWidth sets the width of the frame.
func WithWidth(width int) SetupOption {
	return func(o *SetupOptions) {
		o.Width = width
	}
}

// WithPort sets the port to use for the setup.
func WithPort(port int) SetupOption {
	return func(o *SetupOptions) {
		o.TTY.Port = port
	}
}

// WithFontSize sets the font size for the setup.
func WithFontSize(size int) SetupOption {
	return func(o *SetupOptions) {
		o.TTY.FontSize = size
	}
}

// WithFontFamily sets the font family for the setup.
func WithFontFamily(family string) SetupOption {
	return func(o *SetupOptions) {
		o.TTY.FontFamily = family
	}
}

// WithLineHeight sets the line height for the setup.
func WithLineHeight(height float64) SetupOption {
	return func(o *SetupOptions) {
		o.TTY.LineHeight = height
	}
}

// WithOutput sets the output file for the GIF.
func WithOutput(output string) SetupOption {
	return func(o *SetupOptions) {
		o.Output = output
	}
}

// WithDebug sets the debug flag for setup.
func WithDebug() SetupOption {
	return func(o *SetupOptions) {
		o.TTY.Debug = true
	}
}

// WithPadding sets the padding for the session.
func WithPadding(p string) SetupOption {
	return func(o *SetupOptions) {
		o.Padding = p
	}
}

// New sets up ttyd and go-rod for recording frames.
func New(opts ...SetupOption) Dolly {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	// Get a random port when port is 0.
	if options.TTY.Port == 0 {
		addr, _ := net.Listen("tcp", ":0")
		addr.Close()
		options.TTY.Port = addr.Addr().(*net.TCPAddr).Port
	}
	tty := ttyd.Start(options.TTY)
	go tty.Run()

	os.MkdirAll(options.Folder, os.ModePerm)

	browser := rod.New().MustConnect()

	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", options.TTY.Port))
	page = page.MustSetViewport(options.Width, options.Height, 1, false)
	page = page.MustWaitLoad()
	page = page.MustWaitIdle()
	page.MustElement(".xterm").Eval(fmt.Sprintf(`this.style.padding = '%s'`, options.Padding))
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
				ffmpeg.WithOutput(options.Output),
				ffmpeg.WithWidth(options.Width),
			).Run()

			// Cleanup frames if we successfully made the GIF.
			if err == nil {
				os.RemoveAll(options.Folder)
			}
		},
	}
}
