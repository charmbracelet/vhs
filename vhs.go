package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

// VHS is the object that controls the setup.
type VHS struct {
	Options      *Options
	Page         *rod.Page
	browser      *rod.Browser
	TextCanvas   *rod.Element
	CursorCanvas *rod.Element
	mutex        *sync.Mutex
	recording    bool
	tty          *exec.Cmd
}

// Options is the set of options for the setup.
type Options struct {
	FontFamily    string
	FontSize      int
	LetterSpacing float64
	LineHeight    float64
	Prompt        string
	TypingSpeed   time.Duration
	Theme         Theme
	Test          TestOptions
	Video         VideoOptions
}

const defaultFontSize = 22
const typingSpeed = 50 * time.Millisecond

// DefaultVHSOptions returns the default set of options to use for the setup function.
func DefaultVHSOptions() Options {
	return Options{
		Prompt:        "\\[\\e[38;2;90;86;224m\\]> \\[\\e[0m\\]",
		FontFamily:    "JetBrains Mono,DejaVu Sans Mono,Menlo,Bitstream Vera Sans Mono,Inconsolata,Roboto Mono,Hack,Consolas,ui-monospace,monospace",
		FontSize:      defaultFontSize,
		LetterSpacing: 0,
		LineHeight:    1.0,
		TypingSpeed:   typingSpeed,
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
		recording: true,
		mutex:     mu,
	}
}

// Setup sets up the VHS instance and performs the necessary actions to reflect
// the options that are default and set by the user.
func (vhs *VHS) Setup() {
	// Set Viewport to the correct size, accounting for the padding that will be
	// added during the render.
	padding := vhs.Options.Video.Padding
	width := vhs.Options.Video.Width - padding - padding
	height := vhs.Options.Video.Height - padding - padding
	vhs.Page = vhs.Page.MustSetViewport(width, height, 0, false)

	// Let's wait until we can access the window.term variable.
	vhs.Page = vhs.Page.MustWait("() => window.term != undefined")

	// Find xterm.js canvases for the text and cursor layer for recording.
	vhs.TextCanvas, _ = vhs.Page.Element("canvas.xterm-text-layer")
	vhs.CursorCanvas, _ = vhs.Page.Element("canvas.xterm-cursor-layer")

	// Set Prompt
	vhs.Page.MustElement("textarea").
		MustInput(fmt.Sprintf(` set +o history; export PS1="%s"; clear;`, vhs.Options.Prompt)).
		MustType(input.Enter)

	// Apply options to the terminal
	// By this point the setting commands have been executed, so the `opts` struct is up to date.
	vhs.Page.MustEval(fmt.Sprintf("() => { term.options = { fontSize: %d, fontFamily: '%s', letterSpacing: %f, lineHeight: %f, theme: %s } }",
		vhs.Options.FontSize, vhs.Options.FontFamily, vhs.Options.LetterSpacing,
		vhs.Options.LineHeight, vhs.Options.Theme.String()))

	// Fit the terminal into the window
	vhs.Page.MustEval("term.fit")

	_ = os.RemoveAll(vhs.Options.Video.Input)
	_ = os.MkdirAll(vhs.Options.Video.Input, os.ModePerm)
}

const cleanupWaitTime = 100 * time.Millisecond

// Cleanup cleans up a VHS instance and terminates the go-rod browser and ttyd
// processes.
//
// It also begins the rendering process of the frames into videos.
func (vhs *VHS) Cleanup() {
	vhs.PauseRecording()

	// Give some time for any commands executed (such as `rm`) to finish.
	//
	// If a user runs a long running command, they must sleep for the required time
	// to finish.
	time.Sleep(cleanupWaitTime)

	// Tear down the processes we started.
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
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
		}
	}

	// Cleanup frames if we successfully made the GIF.
	if vhs.Options.Video.CleanupFrames {
		_ = os.RemoveAll(vhs.Options.Video.Input)
	}
}

const quality = 0.92

// Record begins the goroutine which captures images from the xterm.js canvases.
func (vhs *VHS) Record() {
	interval := time.Second / time.Duration(vhs.Options.Video.Framerate)
	time.Sleep(interval)
	go func() {
		counter := 0
		for {
			if !vhs.recording {
				time.Sleep(interval + interval)
				continue
			}
			if vhs.Page != nil {
				counter++
				start := time.Now()
				text, textErr := vhs.TextCanvas.CanvasToImage("image/png", quality)
				cursor, cursorErr := vhs.CursorCanvas.CanvasToImage("image/png", quality)
				if textErr == nil && cursorErr == nil {
					_ = os.WriteFile(vhs.Options.Video.Input+fmt.Sprintf(textFrameFormat, counter), text, os.ModePerm)
					_ = os.WriteFile(vhs.Options.Video.Input+fmt.Sprintf(cursorFrameFormat, counter), cursor, os.ModePerm)
				}
				elapsed := time.Since(start)
				if elapsed >= interval {
					continue
				} else {
					time.Sleep(interval - elapsed)
				}
			}
		}
	}()
}

// ResumeRecording indicates to VHS that the recording should be resumed.
func (vhs *VHS) ResumeRecording() {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	vhs.recording = true
}

// PauseRecording indicates to VHS that the recording should be paused.
func (vhs *VHS) PauseRecording() {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	vhs.recording = false
}
