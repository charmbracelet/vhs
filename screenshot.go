package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	Input string

	style *StyleOptions
}

// NewScreenshotOptions returns ScreenshotOptions by given input.
func NewScreenshotOptions(input string, style *StyleOptions) ScreenshotOptions {
	return ScreenshotOptions{
		frameCapture:       false,
		nextScreenshotPath: "",
		screenshots:        make(map[string]int),
		Input:              input,
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

	fmt.Printf("Total screenshots: %d\n", len(opts.screenshots))

	if len(opts.screenshots) == 0 {
		log.Fatal("NO SCREENSHOTS")
	}

	/*
		TODO
		1 - Allow overwrite
		2 - Check why screenshot sometimes fails
		3 - Refactor styling
	*/

	for path, frame := range opts.screenshots {
		cursorFrame := filepath.Join(opts.Input, fmt.Sprintf(cursorFrameFormat, frame))
		textFrame := filepath.Join(opts.Input, fmt.Sprintf(textFrameFormat, frame))

		fmt.Printf("Making screenshot %s\n", path)

		var args []string
		args = append(args,
			"-y",
			"-i", textFrame,
			"-i", cursorFrame,
		)

		args = append(args, buildFFStyleOpts(opts, path)...)

		fmt.Println(args)

		//nolint:gosec
		cmds = append(cmds, exec.Command(
			"ffmpeg",
			args...,
		))
	}

	return cmds
}

// buildFFopts assembles an ffmpeg command from some VideoOptions
func buildFFStyleOpts(opts ScreenshotOptions, targetFile string) []string {
	// Variables used for building ffmpeg command
	var filterCode strings.Builder
	var args []string
	var prevStageName string
	// stream counter by default 2
	// cursor frame png + text frame png
	streamCounter := 2

	// Compute dimensions of terminal
	termWidth := opts.style.Width
	termHeight := opts.style.Height
	if opts.style.MarginFill != "" {
		termWidth = termWidth - double(opts.style.Margin)
		termHeight = termHeight - double(opts.style.Margin)
	}
	if opts.style.WindowBar != "" {
		termHeight = termHeight - opts.style.WindowBarSize
	}

	// Add margin stream if one is provided
	var marginStream int
	if opts.style.MarginFill != "" {
		if marginFillIsColor(opts.style.MarginFill) {
			// Create plain color stream
			args = append(args,
				"-f", "lavfi",
				"-i",
				fmt.Sprintf(
					"color=%s:s=%dx%d",
					opts.style.MarginFill,
					opts.style.Width,
					opts.style.Height,
				),
			)
			marginStream = streamCounter
			streamCounter++
		} else {
			// Check for existence first.
			_, err := os.Stat(opts.style.MarginFill)
			if err != nil {
				fmt.Println(ErrorStyle.Render("Unable to read margin file: "), opts.style.MarginFill)
			}

			// Add image stream
			args = append(args,
				"-loop", "1",
				"-i", opts.style.MarginFill,
			)
			marginStream = streamCounter
			streamCounter++
		}
	}

	// Create and add a window bar stream if necessary
	var barStream int
	if opts.style.WindowBar != "" {
		barPath := filepath.Join(opts.Input, "bar.png")
		MakeWindowBar(termWidth, termHeight, *opts.style, barPath)
		args = append(args,
			"-i", barPath,
		)
		barStream = streamCounter
		streamCounter++
	}

	// Create and add rounded-corner mask if necessary
	var cornerMaskStream int
	if opts.style.BorderRadius != 0 {
		borderMaskPath := filepath.Join(opts.Input, "mask.png")
		if opts.style.WindowBar != "" {
			MakeBorderRadiusMask(termWidth, termHeight+opts.style.WindowBarSize, opts.style.BorderRadius, borderMaskPath)
		} else {
			MakeBorderRadiusMask(termWidth, termHeight, opts.style.BorderRadius, borderMaskPath)
		}

		args = append(args,
			"-i", borderMaskPath,
		)
		cornerMaskStream = streamCounter
		streamCounter++
	}

	// The following filters are always used
	filterCode.WriteString(
		fmt.Sprintf(`
		[0][1]overlay[merged];
		[merged]scale=%d:%d:force_original_aspect_ratio=1[scaled];
		[scaled]pad=%d:%d:(ow-iw)/2:(oh-ih)/2:%s[padded];
		[padded]fillborders=left=%d:right=%d:top=%d:bottom=%d:mode=fixed:color=%s[padded]
		`,
			termWidth-double(opts.style.Padding),
			termHeight-double(opts.style.Padding),

			termWidth,
			termHeight,
			opts.style.BackgroundColor,

			opts.style.Padding,
			opts.style.Padding,
			opts.style.Padding,
			opts.style.Padding,
			opts.style.BackgroundColor,
		),
	)
	prevStageName = "padded"
	// Add a bar to the terminal and mask the output.
	// This allows us to round the corners of the terminal.
	if opts.style.WindowBar != "" {
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
			[%d]loop=-1[loopbar];
			[loopbar][%s]overlay=0:%d[withbar]
			`,
				barStream,
				prevStageName,
				opts.style.WindowBarSize,
			),
		)
		prevStageName = "withbar"
	}

	if opts.style.BorderRadius != 0 {
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
				[%d]loop=-1[loopmask];
				[%s][loopmask]alphamerge[rounded]
				`,
				cornerMaskStream,
				prevStageName,
			),
		)
		prevStageName = "rounded"
	}

	// Overlay terminal on margin
	if opts.style.MarginFill != "" {
		// ffmpeg will complain if the final filter ends with a semicolon,
		// so we add one BEFORE we start adding filters.
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
			[%d]scale=%d:%d[bg];
			[bg][%s]overlay=(W-w)/2:(H-h)/2:shortest=1[withbg]
			`,
				marginStream,
				opts.style.Width,
				opts.style.Height,
				prevStageName,
			),
		)
		prevStageName = "withbg"
	}

	args = append(args,
		"-filter_complex", filterCode.String(),
		"-map", "["+prevStageName+"]",
		targetFile,
	)

	return args
}
