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

const defaultFrameFileFormat = "frame-%s-%05d.png"

// randomDir returns a random temporary directory to be used for storing frames
// from screenshots of the terminal.
func randomDir() string {
	tmp, _ := os.MkdirTemp(os.TempDir(), "vhs")
	return tmp + "/"
}

// VideoOutputs is a mapping from file type to file path for all video outputs
// of VHS.
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
	Input:         randomDir(),
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

	return exec.Command(
		"ffmpeg", "-y",
		"-framerate", fmt.Sprint(opts.Framerate),
		"-i", opts.Input+"frame-text-%05d.png",
		"-framerate", fmt.Sprint(opts.Framerate),
		"-i", opts.Input+"frame-cursor-%05d.png",
		"-filter_complex",
		"[0][1]overlay[merged];[merged]scale=1000:-1[scaled];[scaled]split[s0][s1];[s0]palettegen=max_colors=256[p];[s1][p]paletteuse[out]",
		"-map", "[out]",
		opts.Output.GIF,
	)
}

// MakeWebM takes a list of images (as frames) and converts them to a WebM.
func MakeWebM(opts VideoOptions) *exec.Cmd {
	if opts.Output.WebM == "" {
		return nil
	}

	fmt.Println("Creating WebM...")

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

	fmt.Println("Creating MP4...")

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
