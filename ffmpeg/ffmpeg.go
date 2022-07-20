package ffmpeg

import (
	"fmt"
	"os/exec"
)

// Options is the set of options for converting frames to a GIF.
type Options struct {
	Input     string
	Output    string
	Framerate int
	Width     int
	MaxColors int
}

// DefaultOptions returns the default set of options for ffmpeg
func DefaultOptions() Options {
	return Options{
		Width:     1200,
		Input:     "tmp/frame-%02d.png",
		Output:    "out.gif",
		Framerate: 50,
		MaxColors: 256,
	}
}

// Option is a function that can be used to set options.
type Option func(*Options)

// WithWidth sets the width of the frame.
func WithWidth(width int) Option {
	return func(o *Options) {
		o.Width = width
	}
}

// WithInput sets the input path for the ffmpeg command.
func WithInput(input string) Option {
	return func(o *Options) {
		o.Input = input
	}
}

// WithOutput sets the output path for the ffmpeg command.
func WithOutput(output string) Option {
	return func(o *Options) {
		o.Output = output
	}
}

// WithFramerate sets the framerate of the GIF.
func WithFramerate(fps int) Option {
	return func(o *Options) {
		o.Framerate = fps
	}
}

// WithMaxColors sets the maximum number of colors for the GIF.
func WithMaxColors(maxColors int) Option {
	return func(o *Options) {
		o.MaxColors = maxColors
	}
}

// MakeGIF takes a list of images (as frames) and converts them to a GIF.
func MakeGIF(opts ...Option) *exec.Cmd {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	flags := fmt.Sprintf(
		"fps=%d,scale=%d:-1:flags=%s,split[s0][s1];[s0]palettegen=max_colors=%d[p];[s1][p]paletteuse",
		options.Framerate,
		options.Width,
		"lanczos",
		options.MaxColors,
	)
	return exec.Command(
		"ffmpeg", "-y", "-i", options.Input,
		"-framerate", fmt.Sprint(options.Framerate),
		"-vf", flags, options.Output,
	)
}
