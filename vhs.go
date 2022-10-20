package vhs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// VHS is the object that controls the setup.
type VHS struct {
	Options   *VHSOptions
	Page      *rod.Page
	browser   *rod.Browser
	mutex     *sync.Mutex
	recording bool
	tty       *exec.Cmd
}

// VHSOptions is the set of options for the setup.
type VHSOptions struct {
	Framerate     float64
	Height        int
	Padding       string
	Prompt        string
	Width         int
	FontFamily    string
	FontSize      int
	LetterSpacing float64
	LineHeight    float64
	TypingSpeed   time.Duration
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
		Prompt:        "\\[\\e[38;2;90;86;224m\\]> \\[\\e[0m\\]",
		FontFamily:    "JetBrains Mono,DejaVu Sans Mono,Menlo,Bitstream Vera Sans Mono,Inconsolata,Roboto Mono,Hack,Consolas,ui-monospace,monospace",
		FontSize:      22,
		LetterSpacing: 0,
		LineHeight:    1.0,
		TypingSpeed:   50 * time.Millisecond,
		Theme:         DefaultTheme,
		Video:         DefaultVideoOptions,
	}
}

// New sets up ttyd and go-rod for recording frames.
func New() VHS {
	port := randomPort()
	tty := StartTTY(port)
	go tty.Run() //nolint:errcheck

	opts := DefaultVHSOptions()
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect().SlowMotion(opts.TypingSpeed)
	page := browser.MustPage(fmt.Sprintf("http://localhost:%d", port))

	mu := &sync.Mutex{}

	return VHS{
		Options:   &opts,
		Page:      page,
		browser:   browser,
		tty:       tty,
		recording: false,
		mutex:     mu,
	}
}

func (vhs *VHS) Setup() {
	vhs.Page = vhs.Page.MustSetViewport(vhs.Options.Width, vhs.Options.Height, 0, false)

	// Let's wait until we can access the window.term variable
	vhs.Page = vhs.Page.MustWait("() => window.term != undefined")

	// Set Prompt
	vhs.Page.MustElement("textarea").
		MustInput(fmt.Sprintf(` set +o history; export PS1="%s"; clear;`, vhs.Options.Prompt)).
		MustType(input.Enter)

	// Apply options to the terminal
	// By this point the setting commands have been executed, so the `opts` struct is up to date.
	vhs.Page.MustEval(fmt.Sprintf("() => { term.options = { fontSize: %d, fontFamily: '%s', letterSpacing: %f, lineHeight: %f, theme: %s } }",
		vhs.Options.FontSize, vhs.Options.FontFamily, vhs.Options.LetterSpacing,
		vhs.Options.LineHeight, vhs.Options.Theme.String()))

	vhs.Page.MustElement(".xterm").MustEval(fmt.Sprintf("() => this.style.padding = '%s'", vhs.Options.Padding))
	vhs.Page.MustElement("body").MustEval("() => this.style.overflow = 'hidden'")
	vhs.Page.MustElement("#terminal-container").MustEval("() => this.style.overflow = 'hidden'")
	vhs.Page.MustElement(".xterm-viewport").MustEval("() => this.style.overflow = 'hidden'")

	// Fit the terminal into the window
	vhs.Page.MustEval("term.fit")

	_ = os.MkdirAll(filepath.Dir(vhs.Options.Video.Input), os.ModePerm)
}

func (vhs *VHS) Cleanup() {
	// Tear down the processes we started.
	vhs.PauseRecording()
	vhs.browser.MustClose()
	_ = vhs.tty.Process.Kill()

	// Generate the video(s) with the frames.
	var cmds []*exec.Cmd
	cmds = append(cmds, MakeGIF(vhs.Options.Video))
	cmds = append(cmds, MakeMP4(vhs.Options.Video))
	cmds = append(cmds, MakeWebM(vhs.Options.Video))

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		_ = cmd.Run()
	}

	// Cleanup frames if we successfully made the GIF.
	if vhs.Options.Video.CleanupFrames {
		os.RemoveAll(vhs.Options.Video.Input)
	}
}

func (vhs *VHS) Record() {
	vhs.ResumeRecording()
	go func() {
		counter := 0
		for {
			if !vhs.recording {
				continue
			}
			counter++
			fmt.Printf("%d ", counter)
			if vhs.Page != nil {
				screenshot, err := vhs.Page.Screenshot(false, &proto.PageCaptureScreenshot{})
				if err == nil {
					_ = os.WriteFile(fmt.Sprintf(vhs.Options.Video.Input, counter), screenshot, 0644)
				}
			}
		}
	}()
}

func (vhs *VHS) ResumeRecording() {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	vhs.recording = true
}

func (vhs *VHS) PauseRecording() {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	vhs.recording = false
}
