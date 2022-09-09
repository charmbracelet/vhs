// gif spawns the ffmpeg process to convert the frames, collected by go-rod's
// screenshots into the input folder, to a GIF.
//
// MakeGIF takes several options to modify the behaviour of the ffmpeg process,
// which can be configured through the Set command.
//
// Set MaxColors 256
// Set Output demo.gif
package vhs

import (
	"fmt"
	"os"
	"os/exec"
)

const defaultFrameFileFormat = "frame-%02d.png"

// randomDir returns a random temporary directory to be used for storing frames
// from screenshots of the terminal.
func randomDir() string {
	tmp, _ := os.MkdirTemp(os.TempDir(), "vhs")
	return tmp
}

// Options is the set of options for converting frames to a GIF.
type GIFOptions struct {
	CleanupFrames bool
	Framerate     int
	Input         string
	MaxColors     int
	Output        string
	Width         int
}

// DefaultGIFOptions is the set of default options for converting frames
// to a GIF, which are used if they are not overridden.
var DefaultGIFOptions = GIFOptions{
	CleanupFrames: true,
	Framerate:     50,
	Input:         randomDir() + "/" + defaultFrameFileFormat,
	MaxColors:     256,
	Output:        "out.gif",
	Width:         1200,
}

// MakeGIF takes a list of images (as frames) and converts them to a GIF.
func MakeGIF(opts GIFOptions) *exec.Cmd {
	flags := fmt.Sprintf(
		"fps=%d,scale=%d:-1:flags=%s,split[s0][s1];[s0]palettegen=max_colors=%d[p];[s1][p]paletteuse",
		opts.Framerate,
		opts.Width,
		"lanczos",
		opts.MaxColors,
	)
	return exec.Command(
		"ffmpeg", "-y", "-i", opts.Input,
		"-framerate", fmt.Sprint(opts.Framerate),
		"-vf", flags, opts.Output,
	)
}
