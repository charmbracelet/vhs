package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// VHS is the object that controls the setup.
type VHS struct {
	Options     *Options
	Errors      []error
	mutex       *sync.Mutex
	started     bool
	recording   bool
	executing   bool
	totalFrames int
	currentTerm *Terminal
	mainTerm    *Terminal
	hiddenTerm  *Terminal
}

// Terminal is the ttyd terminal where commands are executed.
type Terminal struct {
	Page         *rod.Page
	browser      *rod.Browser
	textCanvas   *rod.Element
	cursorCanvas *rod.Element
	tty          *exec.Cmd
	close        func() error
}

// New returns new instance of Terminal.
func NewTerminal(shell Shell) (*Terminal, error) {
	port := randomPort()

	tty := buildTtyCmd(port, shell)
	if err := tty.Start(); err != nil {
		return nil, fmt.Errorf("could not start tty: %w", err)
	}

	path, _ := launcher.LookPath()
	enableNoSandbox := os.Getenv("VHS_NO_SANDBOX") != ""
	u, err := launcher.New().Leakless(false).Bin(path).NoSandbox(enableNoSandbox).Launch()
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}
	browser := rod.New().ControlURL(u).MustConnect()
	page, err := browser.Page(proto.TargetCreateTarget{URL: fmt.Sprintf("http://localhost:%d", port)})
	if err != nil {
		return nil, fmt.Errorf("could not open ttyd: %w", err)
	}

	t := &Terminal{
		browser: browser,
		Page:    page,
		tty:     tty,
		close:   browser.Close,
	}

	return t, nil
}

// Options is the set of options for the setup.
type Options struct {
	Shell         Shell
	FontFamily    string
	FontSize      int
	LetterSpacing float64
	LineHeight    float64
	TypingSpeed   time.Duration
	Theme         Theme
	Test          TestOptions
	Video         VideoOptions
	LoopOffset    float64
}

const (
	defaultFontSize      = 22
	defaultTypingSpeed   = 50 * time.Millisecond
	defaultLineHeight    = 1.0
	defaultLetterSpacing = 1.0
	fontsSeparator       = ","
)

var defaultFontFamily = strings.Join([]string{
	"JetBrains Mono",
	"DejaVu Sans Mono",
	"Menlo",
	"Bitstream Vera Sans Mono",
	"Inconsolata",
	"Roboto Mono",
	"Hack",
	"Consolas",
	"ui-monospace",
	"monospace",
}, fontsSeparator)

// DefaultVHSOptions returns the default set of options to use for the setup function.
func DefaultVHSOptions() Options {
	return Options{
		FontFamily:    defaultFontFamily,
		FontSize:      defaultFontSize,
		LetterSpacing: defaultLetterSpacing,
		LineHeight:    defaultLineHeight,
		TypingSpeed:   defaultTypingSpeed,
		Shell:         Shells[defaultShell],
		Theme:         DefaultTheme,
		Video:         DefaultVideoOptions(),
	}
}

// New sets up ttyd and go-rod for recording frames.
func New() VHS {
	mu := &sync.Mutex{}
	opts := DefaultVHSOptions()
	return VHS{
		Options:   &opts,
		recording: true,
		executing: true,
		mutex:     mu,
	}
}

// Start starts ttyd, browser and everything else needed to create the gif.
func (vhs *VHS) Start() error {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	if vhs.started {
		return fmt.Errorf("vhs is already started")
	}

	// Initialice mainTerm and set it as currentTerm
	t, err := NewTerminal(vhs.Options.Shell)
	if err != nil {
		return err
	}

	vhs.mainTerm, vhs.currentTerm = t, t

	return nil
}

// Setup sets up the VHS instance and performs the necessary actions to reflect
// the options that are default and set by the user.
func (vhs *VHS) Setup() {
	// Set Viewport to the correct size, accounting for the padding that will be
	// added during the render.
	padding := vhs.Options.Video.Padding
	margin := 0
	if vhs.Options.Video.MarginFill != "" {
		margin = vhs.Options.Video.Margin
	}
	bar := 0
	if vhs.Options.Video.WindowBar != "" {
		bar = vhs.Options.Video.WindowBarSize
	}
	width := vhs.Options.Video.Width - double(padding) - double(margin)
	height := vhs.Options.Video.Height - double(padding) - double(margin) - bar
	vhs.mainTerm.Page = vhs.mainTerm.Page.MustSetViewport(width, height, 0, false)

	// Let's wait until we can access the window.term variable.
	vhs.mainTerm.Page = vhs.mainTerm.Page.MustWait("() => window.term != undefined")

	// Find xterm.js canvases for the text and cursor layer for recording.
	vhs.mainTerm.textCanvas, _ = vhs.mainTerm.Page.Element("canvas.xterm-text-layer")
	vhs.mainTerm.cursorCanvas, _ = vhs.mainTerm.Page.Element("canvas.xterm-cursor-layer")

	// Apply options to the terminal
	// By this point the setting commands have been executed, so the `opts` struct is up to date.
	vhs.mainTerm.Page.MustEval(fmt.Sprintf("() => { term.options = { fontSize: %d, fontFamily: '%s', letterSpacing: %f, lineHeight: %f, theme: %s } }",
		vhs.Options.FontSize, vhs.Options.FontFamily, vhs.Options.LetterSpacing,
		vhs.Options.LineHeight, vhs.Options.Theme.String()))

	// Fit the terminal into the window
	vhs.mainTerm.Page.MustEval("term.fit")

	_ = os.RemoveAll(vhs.Options.Video.Input)
	_ = os.MkdirAll(vhs.Options.Video.Input, os.ModePerm)
}

const cleanupWaitTime = 100 * time.Millisecond

// Terminate cleans up a VHS instance and terminates the go-rod browser and ttyd
// processes.
func (vhs *VHS) terminate() error {
	// Give some time for any commands executed (such as `rm`) to finish.
	//
	// If a user runs a long running command, they must sleep for the required time
	// to finish.
	time.Sleep(cleanupWaitTime)

	// Tear down the processes we started.
	vhs.mainTerm.browser.MustClose()
	vhs.hiddenTerm.browser.MustClose()

	err := vhs.mainTerm.tty.Process.Kill()
	if err != nil {
		return err
	}

	return vhs.hiddenTerm.tty.Process.Kill()
}

// Cleanup individual frames.
func (vhs *VHS) Cleanup() error {
	return os.RemoveAll(vhs.Options.Video.Input)
}

// Render starts rendering the individual frames into a video.
func (vhs *VHS) Render() error {
	// Apply Loop Offset by modifying frame sequence
	if err := vhs.ApplyLoopOffset(); err != nil {
		return err
	}

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
			log.Println(string(out))
		}
	}

	return nil
}

// ApplyLoopOffset by modifying frame sequence
func (vhs *VHS) ApplyLoopOffset() error {
	if vhs.totalFrames <= 0 {
		return errors.New("no frames")
	}

	loopOffsetPercentage := vhs.Options.LoopOffset

	// Calculate # of frames to offset from LoopOffset percentage
	loopOffsetFrames := int(math.Ceil(loopOffsetPercentage / 100.0 * float64(vhs.totalFrames)))

	// Take care of overflow and keep track of exact offsetPercentage
	loopOffsetFrames = loopOffsetFrames % vhs.totalFrames

	// No operation if nothing to offset
	if loopOffsetFrames <= 0 {
		return nil
	}

	// Move all frames in [offsetStart, offsetEnd] to end of frame sequence
	offsetStart := vhs.Options.Video.StartingFrame
	offsetEnd := loopOffsetFrames

	// New starting frame will be the next frame after offsetEnd
	vhs.Options.Video.StartingFrame = offsetEnd + 1

	// Rename all text and cursor frame files in the range concurrently
	errCh := make(chan error)
	doneCh := make(chan bool)
	var wg sync.WaitGroup

	for counter := offsetStart; counter <= offsetEnd; counter++ {
		wg.Add(1)
		go func(frameNum int) {
			defer wg.Done()
			offsetFrameNum := frameNum + vhs.totalFrames
			if err := os.Rename(
				filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(cursorFrameFormat, frameNum)),
				filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(cursorFrameFormat, offsetFrameNum)),
			); err != nil {
				errCh <- fmt.Errorf("error applying offset to cursor frame: %w", err)
			}
		}(counter)

		wg.Add(1)
		go func(frameNum int) {
			defer wg.Done()
			offsetFrameNum := frameNum + vhs.totalFrames
			if err := os.Rename(
				filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(textFrameFormat, frameNum)),
				filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(textFrameFormat, offsetFrameNum)),
			); err != nil {
				errCh <- fmt.Errorf("error applying offset to text frame: %w", err)
			}
		}(counter)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		// Bail out in case of an error while renaming
		return err
	}
}

const quality = 1.0

// Record begins the goroutine which captures images from the xterm.js canvases.
func (vhs *VHS) Record(ctx context.Context) <-chan error {
	ch := make(chan error)
	interval := time.Second / time.Duration(vhs.Options.Video.Framerate)

	go func() {
		counter := 0
		start := time.Now()
		for {
			select {
			case <-ctx.Done():
				_ = vhs.terminate()

				// Save total # of frames for offset calculation
				vhs.totalFrames = counter

				// Signal caller that we're done recording.
				close(ch)
				return

			case <-time.After(interval - time.Since(start)):
				// record last attempt
				start = time.Now()

				if !vhs.recording {
					continue
				}
				if vhs.mainTerm.Page == nil {
					continue
				}

				cursor, cursorErr := vhs.mainTerm.cursorCanvas.CanvasToImage("image/png", quality)
				text, textErr := vhs.mainTerm.textCanvas.CanvasToImage("image/png", quality)
				if textErr != nil || cursorErr != nil {
					ch <- fmt.Errorf("error: %v, %v", textErr, cursorErr)
					continue
				}

				counter++

				if err := os.WriteFile(
					filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(cursorFrameFormat, counter)),
					cursor,
					os.ModePerm,
				); err != nil {
					ch <- fmt.Errorf("error writing cursor frame: %w", err)
					continue
				}
				if err := os.WriteFile(
					filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(textFrameFormat, counter)),
					text,
					os.ModePerm,
				); err != nil {
					ch <- fmt.Errorf("error writing text frame: %w", err)
					continue
				}
			}
		}
	}()

	return ch
}

func (vhs *VHS) Close() {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	vhs.mainTerm.close()
	vhs.hiddenTerm.close()
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

// ResumeExecuting indicates to VHS that the executing should be resumed.
// When called mainTerminal will be setted as main terminal.
func (vhs *VHS) ResumeExecuting() {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	vhs.currentTerm = vhs.mainTerm
	vhs.executing = true
}

// PauseExecuting indicates to VHS that the executing should be paused.
// When called hiddenTerm will be setted as currentTerm.
// When executing = false, commands are executed into hidden terminal
// in order to avoid adding those frames into the output.
func (vhs *VHS) PauseExecuting() {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	// If hidden term initialice it.
	// It need page, textCanvas and cursorCanvas in order to execute all commands
	// in hidden ttyd terminal.
	if vhs.hiddenTerm == nil {
		vhs.hiddenTerm, _ = NewTerminal(vhs.Options.Shell)

		vhs.hiddenTerm.Page = vhs.hiddenTerm.Page.MustWait("() => window.term != undefined")
		vhs.hiddenTerm.textCanvas, _ = vhs.hiddenTerm.Page.Element("canvas.xterm-text-layer")
		vhs.hiddenTerm.cursorCanvas, _ = vhs.hiddenTerm.Page.Element("canvas.xterm-cursor-layer")
	}

	vhs.currentTerm = vhs.hiddenTerm
	vhs.executing = false
}
