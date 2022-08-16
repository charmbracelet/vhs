package dolly

import (
	"fmt"
	"os/exec"
)

const frameFileFormat = "frame-%02d.png"

// Options is the set of options for converting frames to a GIF.
type GIFOptions struct {
	InputFolder string
	Output      string
	Framerate   int
	Width       int
	MaxColors   int
}

var DefaultGIFOptions = GIFOptions{
	Width:       1200,
	InputFolder: randomDir(),
	Output:      "out.gif",
	Framerate:   50,
	MaxColors:   256,
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
		"ffmpeg", "-y", "-i", opts.InputFolder+"/"+frameFileFormat,
		"-framerate", fmt.Sprint(opts.Framerate),
		"-vf", flags, opts.Output,
	)
}
