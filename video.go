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
)

const textFrameFormat = "frame-text-%05d.png"
const cursorFrameFormat = "frame-cursor-%05d.png"

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
	GIF  string
	WebM string
	MP4  string
}

// Options is the set of options for converting frames to a GIF.
type VideoOptions struct {
	CleanupFrames   bool
	Framerate       int
	PlaybackSpeed   float64
	Input           string
	MaxColors       int
	Output          VideoOutputs
	Width           int
	Height          int
	Padding         int
	BackgroundColor string
}

const defaultFramerate = 50
const defaultHeight = 600
const defaultMaxColors = 256
const defaultPadding = 72
const defaultPlaybackSpeed = 1.0
const defaultWidth = 1200

// DefaultVideoOptions is the set of default options for converting frames
// to a GIF, which are used if they are not overridden.
func DefaultVideoOptions() VideoOptions {
	return VideoOptions{
		CleanupFrames:   true,
		Framerate:       defaultFramerate,
		Input:           randomDir(),
		MaxColors:       defaultMaxColors,
		Output:          VideoOutputs{GIF: "out.gif", WebM: "", MP4: ""},
		Width:           defaultWidth,
		Height:          defaultHeight,
		Padding:         defaultPadding,
		PlaybackSpeed:   defaultPlaybackSpeed,
		BackgroundColor: DefaultTheme.Background,
	}
}

// MakeGIF takes a list of images (as frames) and converts them to a GIF.
func MakeGIF(opts VideoOptions) *exec.Cmd {
	if opts.Output.GIF == "" {
		return nil
	}

	fmt.Println("Creating GIF...")

	//nolint:gosec
	return exec.Command(
		"ffmpeg", "-y",
		"-r", fmt.Sprint(opts.Framerate),
		"-i", filepath.Join(opts.Input, textFrameFormat),
		"-r", fmt.Sprint(opts.Framerate),
		"-i", filepath.Join(opts.Input, cursorFrameFormat),
		"-filter_complex",
		fmt.Sprintf(`[0][1]overlay[merged];[merged]scale=%d:%d:force_original_aspect_ratio=1[scaled];[scaled]fps=%d,setpts=PTS/%f[speed];[speed]pad=%d:%d:(ow-iw)/2:(oh-ih)/2:%s[padded];[padded]fillborders=left=%d:right=%d:top=%d:bottom=%d:mode=fixed:color=%s[bordered];[bordered]split[a][b];[a]palettegen=max_colors=256[p];[b][p]paletteuse[out]`,
			opts.Width-2*opts.Padding, opts.Height-2*opts.Padding,
			opts.Framerate, opts.PlaybackSpeed,
			opts.Width, opts.Height,
			opts.BackgroundColor,
			opts.Padding, opts.Padding, opts.Padding, opts.Padding,
			opts.BackgroundColor,
		),
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

	//nolint:gosec
	return exec.Command(
		"ffmpeg", "-y",
		"-r", fmt.Sprint(opts.Framerate),
		"-i", filepath.Join(opts.Input, textFrameFormat),
		"-r", fmt.Sprint(opts.Framerate),
		"-i", filepath.Join(opts.Input, cursorFrameFormat),
		"-filter_complex",
		fmt.Sprintf(`[0][1]overlay,scale=%d:%d:force_original_aspect_ratio=1,fps=%d,setpts=PTS/%f,pad=%d:%d:(ow-iw)/2:(oh-ih)/2:%s,fillborders=left=%d:right=%d:top=%d:bottom=%d:mode=fixed:color=%s`,
			opts.Width-2*opts.Padding, opts.Height-2*opts.Padding,
			opts.Framerate, opts.PlaybackSpeed,
			opts.Width, opts.Height,
			opts.BackgroundColor,
			opts.Padding, opts.Padding, opts.Padding, opts.Padding,
			opts.BackgroundColor,
		),
		"-pix_fmt", "yuv420p",
		"-an",
		"-crf", "30",
		"-b:v", "0",
		opts.Output.WebM,
	)
}

// MakeMP4 takes a list of images (as frames) and converts them to an MP4.
func MakeMP4(opts VideoOptions) *exec.Cmd {
	if opts.Output.MP4 == "" {
		return nil
	}

	fmt.Println("Creating MP4...")

	//nolint:gosec
	return exec.Command(
		"ffmpeg", "-y",
		"-r", fmt.Sprint(opts.Framerate),
		"-i", filepath.Join(opts.Input, textFrameFormat),
		"-r", fmt.Sprint(opts.Framerate),
		"-i", filepath.Join(opts.Input, cursorFrameFormat),
		"-filter_complex",
		fmt.Sprintf(`[0][1]overlay,scale=%d:%d:force_original_aspect_ratio=1,fps=%d,setpts=PTS/%f,pad=%d:%d:(ow-iw)/2:(oh-ih)/2:%s,fillborders=left=%d:right=%d:top=%d:bottom=%d:mode=fixed:color=%s`,
			opts.Width-2*opts.Padding, opts.Height-2*opts.Padding,
			opts.Framerate, opts.PlaybackSpeed,
			opts.Width, opts.Height,
			opts.BackgroundColor,
			opts.Padding, opts.Padding, opts.Padding, opts.Padding,
			opts.BackgroundColor,
		),
		"-vcodec", "libx264",
		"-pix_fmt", "yuv420p",
		"-an",
		"-crf", "20",
		opts.Output.MP4,
	)
}
