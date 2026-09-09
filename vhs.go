package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// VHS is the object that controls the setup.
type VHS struct {
	Options        *Options
	Errors         []error
	Page           *rod.Page
	browser        *rod.Browser
	TextCanvas     *rod.Element
	CursorCanvas   *rod.Element
	mutex          *sync.Mutex
	started        bool
	recording      bool
	tty            *exec.Cmd
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

	// browserStartTimeout is how long to wait for Chrome to start and expose its
	// DevTools endpoint before giving up.
	browserStartTimeout = 30 * time.Second

	// browserPollInterval is how often the DevTools HTTP endpoint is probed
	// while waiting for the browser to start.
	browserPollInterval = 100 * time.Millisecond

	// browserCloseTimeout is how long to wait for Chrome to shut down cleanly
	// before it gets killed.
	browserCloseTimeout = 5 * time.Second
)

var defaultWaitPattern = regexp.MustCompile(">$")

var defaultFontFamily = withSymbolsFallback(strings.Join([]string{
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
}, fontsSeparator))

var symbolsFallback = []string{
	"Apple Symbols",
}

func withSymbolsFallback(font string) string {
	return font + fontsSeparator + strings.Join(symbolsFallback, fontsSeparator)
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

// New sets up ttyd and go-rod for recording frames.
func New() VHS {
	mu := &sync.Mutex{}
	opts := DefaultVHSOptions()
	return VHS{
		Options:   &opts,
		recording: true,
		mutex:     mu,
	}
}

// startBrowser launches a headless Chrome and returns a connected go-rod
// browser.
//
// It does not rely on the "DevTools listening on ws://..." line that Chrome
// prints to stderr, as recent Chrome versions do not reliably print it
// (see https://github.com/charmbracelet/vhs/issues/754), which caused VHS to
// hang forever waiting for it. Instead, Chrome is started on a reserved
// port and the DevTools HTTP endpoint is polled until it responds.
func startBrowser(ctx context.Context) (*rod.Browser, func() error, error) {
	binPath, _ := launcher.LookPath()
	if binPath == "" {
		// No browser found on the system, fall back to downloading one.
		var err error
		binPath, err = launcher.NewBrowser().Get()
		if err != nil {
			return nil, nil, fmt.Errorf("could not find or download a browser: %w", err)
		}
	}

	debugPort := randomPort()

	userDataDir, err := os.MkdirTemp("", "vhs-")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create temporary user data directory: %w", err)
	}

	l := launcher.New().
		Leakless(false).
		Bin(binPath).
		UserDataDir(userDataDir).
		RemoteDebuggingPort(debugPort).
		NoSandbox(os.Getenv("VHS_NO_SANDBOX") != "")

	// binPath is resolved by go-rod's launcher (system lookup or its own
	// download cache) and is not user-controlled input.
	cmd := exec.CommandContext(ctx, binPath, l.FormatArgs()...) //nolint:gosec
	if err := cmd.Start(); err != nil {
		_ = os.RemoveAll(userDataDir)
		return nil, nil, fmt.Errorf("could not start browser: %w", err)
	}

	exited := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(exited)
	}()

	if !waitForDevTools(ctx, debugPort, exited, browserStartTimeout) {
		killBrowser(cmd, userDataDir, exited)
		if ctx.Err() != nil {
			return nil, nil, fmt.Errorf("could not start browser: %w", ctx.Err())
		}
		select {
		case <-exited:
			return nil, nil, errors.New("browser exited unexpectedly before its debugging endpoint was ready")
		default:
			return nil, nil, fmt.Errorf("browser debugging endpoint did not become ready within %s", browserStartTimeout)
		}
	}

	wsURL, err := launcher.ResolveURL(fmt.Sprintf("127.0.0.1:%d", debugPort))
	if err != nil {
		killBrowser(cmd, userDataDir, exited)
		return nil, nil, fmt.Errorf("could not resolve browser debugging endpoint: %w", err)
	}

	browser := rod.New().ControlURL(wsURL)
	if err := browser.Connect(); err != nil {
		killBrowser(cmd, userDataDir, exited)
		return nil, nil, fmt.Errorf("could not connect to browser: %w", err)
	}

	closer := func() error {
		if err := browser.Close(); err != nil {
			// The browser could not be closed gracefully (e.g. it is already
			// gone), make sure the process is terminated.
			_ = cmd.Process.Kill()
		}
		select {
		case <-exited:
		case <-time.After(browserCloseTimeout):
			_ = cmd.Process.Kill()
		}
		_ = os.RemoveAll(userDataDir)
		return nil
	}

	return browser, closer, nil
}

// waitForDevTools polls the DevTools HTTP endpoint until it is ready to accept
// connections. It returns true once the endpoint is ready, and false if the
// context is canceled, the browser exits first, or the timeout elapses.
//
// It is the only reliable way to know the browser is ready, regardless of what
// the browser prints to stdout or stderr.
func waitForDevTools(ctx context.Context, port int, exited <-chan struct{}, timeout time.Duration) bool {
	ticker := time.NewTicker(browserPollInterval)
	defer ticker.Stop()
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-ctx.Done():
			return false
		case <-exited:
			return false
		case <-timeoutCh:
			return false
		case <-ticker.C:
			if probeDevTools(port) == nil {
				return true
			}
		}
	}
}

// probeDevTools checks whether the DevTools HTTP endpoint of a browser
// listening on the given port is ready to accept connections.
func probeDevTools(port int) error {
	probeCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(probeCtx, http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/json/version", port), nil)
	if err != nil {
		return fmt.Errorf("could not probe browser debugging endpoint: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not probe browser debugging endpoint: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from browser debugging endpoint: %d", resp.StatusCode)
	}
	return nil
}

// killBrowser terminates a browser process and removes its user data
// directory.
func killBrowser(cmd *exec.Cmd, userDataDir string, exited <-chan struct{}) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	select {
	case <-exited:
	case <-time.After(browserCloseTimeout):
	}
	_ = os.RemoveAll(userDataDir)
}

// Start starts ttyd, browser and everything else needed to create the gif.
func (vhs *VHS) Start(ctx context.Context) error {
	vhs.mutex.Lock()
	defer vhs.mutex.Unlock()

	if vhs.started {
		return fmt.Errorf("vhs is already started")
	}

	port := randomPort()
	vhs.tty = buildTtyCmd(ctx, port, vhs.Options.Shell)
	if err := vhs.tty.Start(); err != nil {
		return fmt.Errorf("could not start tty: %w", err)
	}

	browser, closeBrowser, err := startBrowser(ctx)
	if err != nil {
		_ = vhs.tty.Process.Kill()
		return fmt.Errorf("could not start browser: %w", err)
	}
	page, err := browser.Page(proto.TargetCreateTarget{URL: fmt.Sprintf("http://localhost:%d", port)})
	if err != nil {
		_ = closeBrowser()
		_ = vhs.tty.Process.Kill()
		return fmt.Errorf("could not open ttyd: %w", err)
	}

	vhs.browser = browser
	vhs.Page = page
	vhs.close = func() error {
		_ = closeBrowser()
		_ = vhs.tty.Process.Kill()
		return nil
	}
	vhs.started = true
	return nil
}

// Setup sets up the VHS instance and performs the necessary actions to reflect
// the options that are default and set by the user.
func (vhs *VHS) Setup() error {
	style := vhs.Options.Video.Style

	// Find xterm.js canvases for the text and cursor layer for recording.
	// These exist as soon as ttyd opens the terminal, so this doesn't need
	// to wait for the viewport/font to be set below.
	vhs.TextCanvas, _ = vhs.Page.Element("canvas.xterm-text-layer")
	vhs.CursorCanvas, _ = vhs.Page.Element("canvas.xterm-cursor-layer")

	// Apply options to the terminal
	// By this point the setting commands have been executed, so the `opts` struct is up to date.
	//
	// This must happen before the viewport is sized below: it determines the
	// pixel size of a rendered character cell, which the Rows/Columns case
	// needs in order to compute the viewport, and FontSize/FontFamily/
	// LetterSpacing/LineHeight all affect it.
	vhs.Page.MustEval(fmt.Sprintf("() => { term.options = { fontSize: %d, fontFamily: '%s', letterSpacing: %f, lineHeight: %f, theme: %s, cursorBlink: %t } }",
		vhs.Options.FontSize, vhs.Options.FontFamily, vhs.Options.LetterSpacing,
		vhs.Options.LineHeight, vhs.Options.Theme.String(), vhs.Options.CursorBlink))

	// Account for the padding, margin, and window bar that will be added
	// during the render to determine the viewport size.
	padding := style.Padding
	margin := 0
	if style.MarginFill != "" {
		margin = style.Margin
	}
	bar := 0
	if style.WindowBar != "" {
		bar = style.WindowBarSize
	}
	if style.Rows > 0 || style.Columns > 0 {
		if err := vhs.resolveRowsColumns(padding, margin, bar); err != nil {
			return err
		}
	}

	width := style.Width - double(padding) - double(margin)
	height := style.Height - double(padding) - double(margin) - bar
	vhs.Page = vhs.Page.MustSetViewport(width, height, 0, false)

	// Fit the terminal into the window
	vhs.Page.MustEval("term.fit")

	_ = os.RemoveAll(vhs.Options.Video.Input)
	_ = os.MkdirAll(vhs.Options.Video.Input, 0o750)

	return nil
}

const (
	dimensionProbeAttempts      = 6
	dimensionCorrectionAttempts = 5
)

// resolveRowsColumns derives the pixel Width/Height needed to render the
// requested number of terminal Rows/Columns, overwriting
// vhs.Options.Video.Style.Width/Height in place. The normal viewport-sizing
// code that follows this call then behaves exactly as it does for a plain
// Set Height/Set Width, since by the time it runs, style.Width/style.Height
// already reflect the resolved grid size.
//
// VHS has no native pty (it drives a headless-browser xterm.js instance via
// ttyd), so there is no Go-side font metrics available. Instead, this probes
// the live terminal at a candidate pixel size, reads back the resulting
// term.cols/term.rows (the same public API xterm's own fit addon exposes),
// and uses that ratio to estimate cell size, then iterates to correct for
// the integer floor-rounding xterm's fit addon applies internally.
//
// Both loops are bounded, and an error is returned rather than settling for a
// grid that doesn't match what was asked for: giving up quietly would render
// the tape at the wrong size, or at 0x0 if the terminal was never measured.
func (vhs *VHS) resolveRowsColumns(padding, margin, bar int) error {
	style := vhs.Options.Video.Style

	measure := func(w, h int) (cols, rows int) {
		vhs.Page = vhs.Page.MustSetViewport(w, h, 0, false)
		vhs.Page.MustEval("term.fit")
		dims := vhs.Page.MustEval("() => ({ cols: term.cols, rows: term.rows })")
		return dims.Get("cols").Int(), dims.Get("rows").Int()
	}

	probeWidth := style.Width - double(padding) - double(margin)
	probeHeight := style.Height - double(padding) - double(margin) - bar

	var cellWidth, cellHeight float64
	var measured bool
	var lastCols, lastRows int
	for range dimensionProbeAttempts {
		cols, rows := measure(probeWidth, probeHeight)
		lastCols, lastRows = cols, rows

		tooSmall := cols <= 0 || rows <= 0 ||
			(style.Columns > 0 && cols < style.Columns) ||
			(style.Rows > 0 && rows < style.Rows)
		if !tooSmall {
			cellWidth = float64(probeWidth) / float64(cols)
			cellHeight = float64(probeHeight) / float64(rows)
			measured = true
			break
		}
		probeWidth = double(probeWidth)
		probeHeight = double(probeHeight)
	}

	// Without a successful measurement there is no cell size to scale by, and
	// continuing would size the viewport to 0x0.
	if !measured {
		return fmt.Errorf(
			"could not fit %s: terminal measured %d columns x %d rows at %dx%d pixels after %d attempts",
			describeGrid(style.Columns, style.Rows),
			lastCols, lastRows, probeWidth, probeHeight, dimensionProbeAttempts,
		)
	}

	contentWidth := probeWidth
	contentHeight := probeHeight
	if style.Columns > 0 {
		contentWidth = int(math.Round(cellWidth * float64(style.Columns)))
	}
	if style.Rows > 0 {
		contentHeight = int(math.Round(cellHeight * float64(style.Rows)))
	}

	// xterm's fit addon floors container-size / cell-size, so the estimate
	// above can be off by a cell. Nudge the content size until the measured
	// cols/rows exactly match what was requested.
	var corrected bool
	for range dimensionCorrectionAttempts {
		cols, rows := measure(contentWidth, contentHeight)
		lastCols, lastRows = cols, rows

		converged := true
		if style.Columns > 0 && cols != style.Columns {
			contentWidth += int(math.Round(float64(style.Columns-cols) * cellWidth))
			converged = false
		}
		if style.Rows > 0 && rows != style.Rows {
			contentHeight += int(math.Round(float64(style.Rows-rows) * cellHeight))
			converged = false
		}
		if converged {
			corrected = true
			break
		}
	}

	// The grid never settled on the requested size, so recording now would
	// silently produce a tape with the wrong number of rows/columns.
	if !corrected {
		return fmt.Errorf(
			"could not resolve a viewport for %s: closest match was %d columns x %d rows after %d attempts",
			describeGrid(style.Columns, style.Rows),
			lastCols, lastRows, dimensionCorrectionAttempts,
		)
	}

	if style.Columns > 0 {
		style.Width = roundUpToEven(contentWidth + double(padding) + double(margin))
	}
	if style.Rows > 0 {
		style.Height = roundUpToEven(contentHeight + double(padding) + double(margin) + bar)
	}

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

// describeGrid names the grid size that was requested, mentioning only the
// dimension(s) the tape actually set.
func describeGrid(columns, rows int) string {
	switch {
	case columns > 0 && rows > 0:
		return fmt.Sprintf("%d columns x %d rows", columns, rows)
	case columns > 0:
		return fmt.Sprintf("%d columns", columns)
	default:
		return fmt.Sprintf("%d rows", rows)
	}
}

const cleanupWaitTime = 100 * time.Millisecond

// Terminate cleans up a VHS instance and terminates the go-rod browser and ttyd
// processes.
//
//nolint:wrapcheck
func (vhs *VHS) terminate() error {
	// Give some time for any commands executed (such as `rm`) to finish.
	//
	// If a user runs a long running command, they must sleep for the required time
	// to finish.
	time.Sleep(cleanupWaitTime)

	// Tear down the processes we started.
	//
	// The browser may have already exited, in which case closing it fails and
	// its process gets killed as a fallback, so errors are ignored here.
	_ = vhs.browser.Close()
	return vhs.tty.Process.Kill()
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

	// Generate the video(s) with the frames.
	var cmds []*exec.Cmd //nolint:prealloc
	cmds = append(cmds, MakeGIF(ctx, vhs.Options.Video))
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

const quality = 1.0

// Record begins the goroutine which captures images from the xterm.js canvases.
func (vhs *VHS) Record(ctx context.Context) <-chan error {
	ch := make(chan error)
	interval := time.Second / time.Duration(vhs.Options.Video.Framerate)

	//nolint: mnd
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
				if vhs.Page == nil {
					continue
				}

				cursor, cursorErr := vhs.CursorCanvas.CanvasToImage("image/png", quality)
				text, textErr := vhs.TextCanvas.CanvasToImage("image/png", quality)
				if textErr != nil || cursorErr != nil {
					ch <- fmt.Errorf("error: %v, %v", textErr, cursorErr)
					continue
				}

				counter++
				if err := os.WriteFile(
					filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(cursorFrameFormat, counter)),
					cursor,
					0o600,
				); err != nil {
					ch <- fmt.Errorf("error writing cursor frame: %w", err)
					continue
				}
				if err := os.WriteFile(
					filepath.Join(vhs.Options.Video.Input, fmt.Sprintf(textFrameFormat, counter)),
					text,
					0o600,
				); err != nil {
					ch <- fmt.Errorf("error writing text frame: %w", err)
					continue
				}

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
