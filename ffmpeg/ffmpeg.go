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

// MakeGIF takes a list of images (as frames) and converts them to a GIF.
func MakeGIF(opts Options) *exec.Cmd {
	flags := fmt.Sprintf(
		"fps=%d,scale=%d:-1:flags=%s,split[s0][s1];[s0]palettegen=max_colors=%d[p];[s1][p]paletteuse",
		opts.Framerate,
		opts.Width,
		"lanczos",
		opts.MaxColors,
	)
	return exec.Command("ffmpeg", "-y", "-i", opts.Input, "-framerate", fmt.Sprint(opts.Framerate), "-vf", flags, opts.Output)
}
