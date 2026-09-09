package main

import (
	"context"
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
	sixel bool
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
func MakeScreenshots(ctx context.Context, opts ScreenshotOptions) []*exec.Cmd {
	cmds := []*exec.Cmd{} //nolint:prealloc

	for path, frame := range opts.screenshots {
		cursorStream := filepath.Join(opts.input, fmt.Sprintf(cursorFrameFormat, frame))
		textStream := filepath.Join(opts.input, fmt.Sprintf(textFrameFormat, frame))
		imageStream := filepath.Join(opts.input, fmt.Sprintf(imageFrameFormat, frame))

		args := opts.buildFFopts(path, textStream, cursorStream, imageStream)

		cmds = append(cmds, exec.CommandContext(
			ctx,
			"ffmpeg",
			args...,
		))
	}

	return cmds
}

// buildFFopts assembles an ffmpeg command from some VideoOptions.
func (opts *ScreenshotOptions) buildFFopts(targetFile, textStream, cursorStream, imageStream string) []string {
	var args []string //nolint:prealloc
	streamCounter := 2
	if opts.sixel {
		streamCounter++
	}

	streamBuilder := NewStreamBuilder(streamCounter, opts.input, opts.style)
	// Input frame options, used no matter what
	// Stream 0: text frames
	// Stream 1: cursor frames
	streamBuilder.args = append(streamBuilder.args,
		"-y",
		"-i", textStream,
		"-i", cursorStream,
	)
	if opts.sixel {
		streamBuilder.args = append(streamBuilder.args, "-i", imageStream)
	}

	streamBuilder = streamBuilder.
		WithMargin().
		WithBar().
		WithCorner()

	filterBuilder := NewScreenshotFilterComplexBuilder(opts.style, opts.sixel).
		WithWindowBar(streamBuilder.barStream).
		WithBorderRadius(streamBuilder.cornerStream).
		WithMarginFill(streamBuilder.marginStream)

	args = append(args, streamBuilder.Build()...)
	args = append(args, filterBuilder.Build()...)
	args = append(args, targetFile)

	return args
}
