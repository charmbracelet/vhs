package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// ScreenshotOptions holds options related with screenshots.
type ScreenshotOptions struct {
	// frameCapture holds a flag indicating if screenshots must be taken.
	frameCapture bool

	// nextScreenshotPath holds the path of new screenshot.
	nextScreenshotPath string

	// screenshots represents a map holding screenshot path as key and frame as value.
	screenshots map[string]int

	// input represents location of cursor and text frames png files.
	input string
}

// NewScreenshotOptions returns ScreenshotOptions by given input.
func NewScreenshotOptions(input string) ScreenshotOptions {
	return ScreenshotOptions{
		frameCapture:       false,
		nextScreenshotPath: "",
		screenshots:        make(map[string]int),
		input:              input,
	}
}

// makeScreenshot stores in screenshots map the target frame of the screenshot.
// After storing frame it disables frame capture.
func (opts *ScreenshotOptions) makeScreenshot(frame int) {
	opts.screenshots[opts.nextScreenshotPath] = frame

	opts.frameCapture = false
	opts.nextScreenshotPath = ""
}

// captureNextFrame prepares capture of next frame by given path.
func (opts *ScreenshotOptions) enableFrameCapture(path string) {
	opts.frameCapture = true
	opts.nextScreenshotPath = path
}

// MakeScreenshots generates screenshots by given ScreenshotOptions.
func MakeScreenshots(opts ScreenshotOptions) []*exec.Cmd {
	cmds := []*exec.Cmd{}

	for path, frame := range opts.screenshots {
		cursorFrame := filepath.Join(opts.input, fmt.Sprintf(cursorFrameFormat, frame))
		textFrame := filepath.Join(opts.input, fmt.Sprintf(textFrameFormat, frame))

		cmds = append(cmds, exec.Command(
			"ffmpeg",
			"-i", textFrame,
			"-i", cursorFrame,
			"-filter_complex",
			"overlay=0:0",
			path,
		))
	}

	return cmds
}
