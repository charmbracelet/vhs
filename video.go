// Package vhs video.go spawns the ffmpeg process to convert the frames,
// collected by go-rod's  screenshots into the input folder, to a GIF, WebM,
// MP4.
//
// MakeGIF takes several options to modify the behaviour of the ffmpeg process,
// which can be configured through the Set command.
//
// Set MaxColors 256
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	textFrameFormat   = "frame-text-%05d.png"
	cursorFrameFormat = "frame-cursor-%05d.png"
)

const (
	mp4  = ".mp4"
	webm = ".webm"
	gif  = ".gif"
)

// randomDir returns a random temporary directory to be used for storing frames
// from screenshots of the terminal.
func randomDir() string {
	tmp, err := os.MkdirTemp(os.TempDir(), "vhs")
	if err != nil {
		log.Printf("Error creating temporary directory: %v", err)
	}
	return tmp
}

// VideoOutputs is a mapping from file type to file path for all video outputs
// of VHS.
type VideoOutputs struct {
	GIF    string
	WebM   string
	MP4    string
	Frames string
}

// VideoOptions is the set of options for converting frames to a GIF.
type VideoOptions struct {
	Framerate     int
	PlaybackSpeed float64
	Input         string
	MaxColors     int
	Output        VideoOutputs
	StartingFrame int
	Style         *StyleOptions
}

const (
	defaultFramerate     = 50
	defaultStartingFrame = 1
)

// DefaultVideoOptions is the set of default options for converting frames
// to a GIF, which are used if they are not overridden.
func DefaultVideoOptions() VideoOptions {
	return VideoOptions{
		Framerate:     defaultFramerate,
		Input:         randomDir(),
		MaxColors:     defaultMaxColors,
		Output:        VideoOutputs{GIF: "", WebM: "", MP4: "", Frames: ""},
		PlaybackSpeed: defaultPlaybackSpeed,
		StartingFrame: defaultStartingFrame,
	}
}

func marginFillIsColor(marginFill string) bool {
	return strings.HasPrefix(marginFill, "#")
}

// ensureDir ensures that the file path of the output can be created by
// creating all the necessary nested folders.
func ensureDir(output string) {
	err := os.MkdirAll(filepath.Dir(output), os.ModePerm)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Unable to create output directory: "), output)
	}
}

// buildFFopts assembles an ffmpeg command from some VideoOptions.
func buildFFopts(opts VideoOptions, targetFile string) []string {
	var args []string
	streamCounter := 2

	streamBuilder := NewStreamBuilder(streamCounter, opts.Input, opts.Style)

	// Input frame options, used no matter what
	// Stream 0: text frames
	// Stream 1: cursor frames
	streamBuilder.args = append(streamBuilder.args,
		"-y",
		"-r", fmt.Sprint(opts.Framerate),
		"-start_number", fmt.Sprint(opts.StartingFrame),
		"-i", filepath.Join(opts.Input, textFrameFormat),
		"-r", fmt.Sprint(opts.Framerate),
		"-start_number", fmt.Sprint(opts.StartingFrame),
		"-i", filepath.Join(opts.Input, cursorFrameFormat),
	)

	streamBuilder = streamBuilder.
		WithMargin().
		WithBar().
		WithCorner()

	filterBuilder := NewVideoFilterBuilder(&opts).
		WithWindowBar(streamBuilder.barStream).
		WithBorderRadius(streamBuilder.cornerStream).
		WithMarginFill(streamBuilder.marginStream)

	// Format-specific options
	switch filepath.Ext(targetFile) {
	case gif:
		filterBuilder = filterBuilder.WithGIF()
	case webm:
		streamBuilder = streamBuilder.WithWebm()
	case mp4:
		streamBuilder = streamBuilder.WithMP4()
	}

	args = append(args, streamBuilder.Build()...)
	args = append(args, filterBuilder.Build()...)
	args = append(args, targetFile)

	return args
}

// MakeGIF takes a list of images (as frames) and converts them to a GIF.
func MakeGIF(opts VideoOptions) *exec.Cmd {
	if opts.Output.GIF == "" {
		return nil
	}

	log.Println(GrayStyle.Render("Creating " + opts.Output.GIF + "..."))
	ensureDir(opts.Output.GIF)

	//nolint:gosec
	return exec.Command(
		"ffmpeg",
		buildFFopts(opts, opts.Output.GIF)...,
	)
}

// MakeWebM takes a list of images (as frames) and converts them to a WebM.
func MakeWebM(opts VideoOptions) *exec.Cmd {
	if opts.Output.WebM == "" {
		return nil
	}

	log.Println(GrayStyle.Render("Creating " + opts.Output.WebM + "..."))
	ensureDir(opts.Output.WebM)

	//nolint:gosec
	return exec.Command(
		"ffmpeg",
		buildFFopts(opts, opts.Output.WebM)...,
	)
}

// MakeMP4 takes a list of images (as frames) and converts them to an MP4.
func MakeMP4(opts VideoOptions) *exec.Cmd {
	if opts.Output.MP4 == "" {
		return nil
	}

	log.Println(GrayStyle.Render("Creating " + opts.Output.MP4 + "..."))
	ensureDir(opts.Output.MP4)

	//nolint:gosec
	return exec.Command(
		"ffmpeg",
		buildFFopts(opts, opts.Output.MP4)...,
	)
}
