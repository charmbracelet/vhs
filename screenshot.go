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

	// Input represents location of cursor and text frames png files.
	input string

	style *StyleOptions
}

// NewScreenshotOptions returns ScreenshotOptions by given input.
func NewScreenshotOptions(input string, style *StyleOptions) ScreenshotOptions {
	return ScreenshotOptions{
		frameCapture:       false,
		nextScreenshotPath: "",
		screenshots:        make(map[string]int),
		input:              input,
		style:              style,
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
		cursorStream := filepath.Join(opts.input, fmt.Sprintf(cursorFrameFormat, frame))
		textStream := filepath.Join(opts.input, fmt.Sprintf(textFrameFormat, frame))

		args := opts.buildFFopts(path, textStream, cursorStream)

		//nolint:gosec
		cmds = append(cmds, exec.Command(
			"ffmpeg",
			args...,
		))
	}

	return cmds
}

// buildFFopts assembles an ffmpeg command from some VideoOptions.
func (opts *ScreenshotOptions) buildFFopts(targetFile, textStream, cursorStream string) []string {
	var args []string
	streamCounter := 2

	streamBuilder := NewStreamBuilder(streamCounter, opts.input, opts.style)
	// Input frame options, used no matter what
	// Stream 0: text frames
	// Stream 1: cursor frames
	streamBuilder.args = append(streamBuilder.args,
		"-y",
		"-i", textStream,
		"-i", cursorStream,
	)

	streamBuilder = streamBuilder.
		WithMargin().
		WithBar().
		WithCorner()

	filterBuilder := NewScreenshotFilterComplexBuilder(opts.style).
		WithWindowBar(streamBuilder.barStream).
		WithBorderRadius(streamBuilder.cornerStream).
		WithMarginFill(streamBuilder.marginStream)

	args = append(args, streamBuilder.Build()...)
	args = append(args, filterBuilder.Build()...)
	args = append(args, targetFile)

	return args
}
