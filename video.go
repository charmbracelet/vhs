// video spawns the ffmpeg process to convert the frames, collected by go-rod's
// screenshots into the input folder, to a GIF, WebM, MP4.
//
// MakeGIF takes several options to modify the behaviour of the ffmpeg process,
// which can be configured through the Set command.
//
// Set MaxColors 256
package vhs

import (
	"fmt"
	"os"
	"os/exec"
)

const defaultFrameFileFormat = "frame-%05d.png"

// randomDir returns a random temporary directory to be used for storing frames
// from screenshots of the terminal.
func randomDir() string {
	tmp, _ := os.MkdirTemp(os.TempDir(), "vhs")
	return tmp
}

type VideoOutputs struct {
	GIF  string
	WebM string
	MP4  string
}

// Options is the set of options for converting frames to a GIF.
type VideoOptions struct {
	CleanupFrames bool
	Framerate     int
	Input         string
	MaxColors     int
	Output        VideoOutputs
	Width         int
}

// DefaultVideoOptions is the set of default options for converting frames
// to a GIF, which are used if they are not overridden.
var DefaultVideoOptions = VideoOptions{
	CleanupFrames: true,
	Framerate:     50,
	Input:         randomDir() + "/" + defaultFrameFileFormat,
	MaxColors:     256,
	Output:        VideoOutputs{GIF: "out.gif", WebM: "", MP4: ""},
	Width:         1200,
}

// MakeGIF takes a list of images (as frames) and converts them to a GIF.
func MakeGIF(opts VideoOptions) *exec.Cmd {
	if opts.Output.GIF == "" {
		return nil
	}

	fmt.Println("Creating GIF...")

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
		"-vf", flags, opts.Output.GIF,
	)
}

// MakeWebM takes a list of images (as frames) and converts them to a WebM.
func MakeWebM(opts VideoOptions) *exec.Cmd {
	if opts.Output.WebM == "" {
		return nil
	}

	return exec.Command(
		"ffmpeg", "-y", "-i", opts.Input,
		"-framerate", fmt.Sprint(opts.Framerate),
		"-pix_fmt", "yuv420p",
		"-an",
		"-crf", "30",
		"-b:v", "0",
		"-filter:v", fmt.Sprintf("scale=%d:-1", opts.Width),
		opts.Output.WebM,
	)
}

// MakeMP4 takes a list of images (as frames) and converts them to an MP4.
func MakeMP4(opts VideoOptions) *exec.Cmd {
	if opts.Output.MP4 == "" {
		return nil
	}

	return exec.Command(
		"ffmpeg", "-y", "-i", opts.Input,
		"-framerate", fmt.Sprint(opts.Framerate),
		"-vcodec", "libx264",
		"-pix_fmt", "yuv420p",
		"-an",
		"-crf", "20",
		"-filter:v", fmt.Sprintf("scale=%d:-1", opts.Width),
		opts.Output.MP4,
	)
}
