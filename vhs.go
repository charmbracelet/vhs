package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// VHS is the object that controls the setup.
type VHS struct {
	Options        *Options
	Errors         []error
	term           *Terminal
	renderer       *Renderer
	mutex          *sync.Mutex
	started        bool
	recording      bool
	totalFrames    int
	close          func() error
	widthExplicit  bool
	heightExplicit bool
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
	WaitTimeout   time.Duration
	WaitPattern   *regexp.Regexp
	CursorBlink   bool
	Screenshot    ScreenshotOptions
	Style         StyleOptions
}

const (
	defaultFontSize      = 22
	defaultTypingSpeed   = 50 * time.Millisecond
	defaultLineHeight    = 1.0
	defaultLetterSpacing = 1.0
	fontsSeparator       = ","
	defaultCursorBlink   = true
	defaultWaitTimeout   = 15 * time.Second

	// The terminal starts at this size, and is resized to fit the video
	// once the settings are known.
	initialColumns = 80
	initialRows    = 24
)

var defaultWaitPattern = regexp.MustCompile(">$")

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

// symbolsFallback are fonts used for symbols the chosen fonts lack.
var symbolsFallback = []string{
	"Apple Symbols",
}

// DefaultVHSOptions returns the default set of options to use for the setup function.
func DefaultVHSOptions() Options {
	style := DefaultStyleOptions()
	video := DefaultVideoOptions()
	video.Style = style
	screenshot := NewScreenshotOptions(video.Input, style)

	return Options{
		FontFamily:    defaultFontFamily,
		FontSize:      defaultFontSize,
		LetterSpacing: defaultLetterSpacing,
		LineHeight:    defaultLineHeight,
		TypingSpeed:   defaultTypingSpeed,
		Shell:         Shells[defaultShell],
		Theme:         DefaultTheme,
		CursorBlink:   defaultCursorBlink,
		Video:         video,
		Screenshot:    screenshot,
		WaitTimeout:   defaultWaitTimeout,
		WaitPattern:   defaultWaitPattern,
	}
}

// New creates a VHS instance with the default options.
func New() VHS {
	mu := &sync.Mutex{}
	opts := DefaultVHSOptions()
	return VHS{
		Options:   &opts,
		recording: true,
		mutex:     mu,
	}
}

// Start starts the shell in a terminal.
func (vhs *VHS) Start(ctx context.Context) error {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	if vhs.started {
		return fmt.Errorf("vhs is already started")
	}

	term, err := StartTerminal(ctx, vhs.Options.Shell, initialColumns, initialRows)
	if err != nil {
		return err
	}

	vhs.term = term
	vhs.close = term.Close
	vhs.started = true
	return nil
}

// Setup sets up the VHS instance and performs the necessary actions to reflect
// the options that are default and set by the user.
func (vhs *VHS) Setup() error {
	style := vhs.Options.Video.Style

	// By this point the setting commands have been executed, so the options
	// are final.
	renderer, err := NewRenderer(vhs.Options)
	if err != nil {
		return err
	}
	vhs.term.SetTheme(vhs.Options.Theme)

	// Account for the padding, margin, and window bar that will be added
	// during the render to determine the terminal size.
	padding := style.Padding
	margin := 0
	if style.MarginFill != "" {
		margin = style.Margin
	}
	bar := 0
	if style.WindowBar != "" {
		bar = style.WindowBarSize
	}

	cellW, cellH := renderer.CellSize()
	if style.Columns > 0 {
		style.Width = roundUpToEven(style.Columns*cellW + double(padding) + double(margin))
	}
	if style.Rows > 0 {
		style.Height = roundUpToEven(style.Rows*cellH + double(padding) + double(margin) + bar)
	}

	width := style.Width - double(padding) - double(margin)
	height := style.Height - double(padding) - double(margin) - bar
	cols, rows := 0, 0
	if width > 0 && height > 0 {
		cols, rows = width/cellW, height/cellH
	}
	if cols < 1 || rows < 1 {
		return fmt.Errorf(
			"could not fit the terminal in %dx%d pixels: each cell is %dx%d pixels",
			width, height, cellW, cellH,
		)
	}

	renderer.SetSize(width, height, cols, rows)
	if err := vhs.term.Resize(cols, rows); err != nil {
		return err
	}
	vhs.renderer = renderer

	_ = os.RemoveAll(vhs.Options.Video.Input)
	_ = os.MkdirAll(vhs.Options.Video.Input, 0o750)

	return nil
}

// roundUpToEven rounds n up to the nearest even number. MP4 and WebM
// encoders require even width/height, and the grid-derived dimensions above
// can land on an odd pixel count.
func roundUpToEven(n int) int {
	if n%2 != 0 {
		return n + 1
	}
	return n
}

const cleanupWaitTime = 100 * time.Millisecond

// terminate cleans up a VHS instance and stops the shell.
func (vhs *VHS) terminate() error {
	// Give some time for any commands executed (such as `rm`) to finish.
	//
	// If a user runs a long running command, they must sleep for the required time
	// to finish.
	time.Sleep(cleanupWaitTime)
	return vhs.term.Close()
}

// Cleanup individual frames.
//
//nolint:wrapcheck
func (vhs *VHS) Cleanup() error {
	err := os.RemoveAll(vhs.Options.Video.Input)
	if err != nil {
		return err
	}
	return os.RemoveAll(vhs.Options.Screenshot.input)
}

// Render starts rendering the individual frames into a video.
func (vhs *VHS) Render(ctx context.Context) error {
	// Apply Loop Offset by modifying frame sequence
	if err := vhs.ApplyLoopOffset(); err != nil {
		return err
	}

	// The frames are on disk: give the memory used to record them back to
	// the system before ffmpeg needs its own.
	debug.FreeOSMemory()

	// Generate the video(s) with the frames.
	var cmds []*exec.Cmd //nolint:prealloc
	cmds = append(cmds, MakeGIF(ctx, vhs.Options.Video)...)
	cmds = append(cmds, MakeMP4(ctx, vhs.Options.Video))
	cmds = append(cmds, MakeWebM(ctx, vhs.Options.Video))
	cmds = append(cmds, MakeScreenshots(ctx, vhs.Options.Screenshot)...)

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(string(out))
			return fmt.Errorf("render failed: %w", err)
		}
	}

	return nil
}

// ApplyLoopOffset by modifying frame sequence.
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

// encodedLayer is a layer of a frame encoded to PNG, shared by the frames in
// which it looks the same.
type encodedLayer struct {
	done chan struct{}
	png  []byte
	err  error
}

// pngBuffers lets the PNG encoder reuse its buffers from frame to frame.
type pngBuffers struct{ pool sync.Pool }

func (p *pngBuffers) Get() *png.EncoderBuffer {
	b, _ := p.pool.Get().(*png.EncoderBuffer)
	return b
}

func (p *pngBuffers) Put(b *png.EncoderBuffer) { p.pool.Put(b) }

var pngEncoder = png.Encoder{CompressionLevel: png.BestSpeed, BufferPool: &pngBuffers{}}

// writeFrame writes the layers of frame n once they are encoded.
func writeFrame(dir string, n int, text, cursor *encodedLayer) error {
	for _, l := range []struct {
		layer  *encodedLayer
		format string
	}{{cursor, cursorFrameFormat}, {text, textFrameFormat}} {
		<-l.layer.done
		if l.layer.err != nil {
			return l.layer.err
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf(l.format, n)), l.layer.png, 0o600); err != nil {
			return fmt.Errorf("error writing frame: %w", err)
		}
	}
	return nil
}

// Record begins the goroutine which renders the terminal into frames.
//
// Frames are encoded and written in the background so that slow encoding
// doesn't skew the frame timing. The text and cursor layers of a frame are
// only encoded when they change: frames that look the same as the previous
// one, in either layer, reuse its encoding.
func (vhs *VHS) Record(ctx context.Context) <-chan error {
	ch := make(chan error)
	interval := time.Second / time.Duration(vhs.Options.Video.Framerate)

	go func() {
		var wg sync.WaitGroup
		// Frames being encoded hold their images: bound how many.
		encoders := make(chan struct{}, max(2, runtime.NumCPU()/2)) //nolint:mnd
		encode := func(img *image.RGBA, release func(*image.RGBA)) *encodedLayer {
			l := &encodedLayer{done: make(chan struct{})}
			encoders <- struct{}{}
			wg.Go(func() {
				defer func() { <-encoders }()
				defer close(l.done)
				var b bytes.Buffer
				if err := pngEncoder.Encode(&b, img); err != nil {
					l.err = fmt.Errorf("error encoding frame: %w", err)
				}
				l.png = b.Bytes()
				release(img)
			})
			return l
		}
		var text, cursor *encodedLayer
		counter := 0
		start := time.Now()
		for {
			select {
			case <-ctx.Done():
				wg.Wait()
				_ = vhs.terminate()

				// Making the video only needs the frames: let go of the
				// terminal and what drew it while ffmpeg runs.
				_ = vhs.renderer.Close()
				vhs.renderer, vhs.term = nil, nil

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

				t, c := vhs.renderer.Frame(vhs.term, start)
				if t != nil {
					text = encode(t, vhs.renderer.ReleaseText)
				}
				if c != nil {
					cursor = encode(c, vhs.renderer.ReleaseCursor)
				}

				counter++
				text, cursor, n := text, cursor, counter
				wg.Go(func() {
					if err := writeFrame(vhs.Options.Video.Input, n, text, cursor); err != nil {
						ch <- err
					}
				})

				// Capture current frame and disable frame capturing
				if vhs.Options.Screenshot.frameCapture {
					vhs.Options.Screenshot.makeScreenshot(counter)
				}
			}
		}
	}()

	return ch
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

// ScreenshotNextFrame indicates to VHS that screenshot of next frame must be taken.
func (vhs *VHS) ScreenshotNextFrame(path string) {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	vhs.Options.Screenshot.enableFrameCapture(path)
}
