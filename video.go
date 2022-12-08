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

// VideoOptions is the set of options for converting frames to a GIF.
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
	StartingFrame   int
	ImageFrame      string
	FrameSize       int
}

const (
	defaultFramerate     = 50
	defaultHeight        = 600
	defaultMaxColors     = 256
	defaultPadding       = 72
	defaultPlaybackSpeed = 1.0
	defaultWidth         = 1200
	defaultStartingFrame = 1
)

// DefaultVideoOptions is the set of default options for converting frames
// to a GIF, which are used if they are not overridden.
func DefaultVideoOptions() VideoOptions {
	return VideoOptions{
		ImageFrame:      "",
		FrameSize:       25,
		CleanupFrames:   true,
		Framerate:       defaultFramerate,
		Input:           randomDir(),
		MaxColors:       defaultMaxColors,
		Output:          VideoOutputs{GIF: "", WebM: "", MP4: ""},
		Width:           defaultWidth,
		Height:          defaultHeight,
		Padding:         defaultPadding,
		PlaybackSpeed:   defaultPlaybackSpeed,
		BackgroundColor: DefaultTheme.Background,
		StartingFrame:   defaultStartingFrame,
	}
}

// ensureDir ensures that the file path of the output can be created by
// creating all the necessary nested folders.
func ensureDir(output string) {
	err := os.MkdirAll(filepath.Dir(output), os.ModePerm)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Unable to create output directory: "), output)
	}
}

// BuildFFopts builds an ffmpeg command from a VideoOptions
func buildFFopts(opts VideoOptions, targetFile string) []string {
	// Input frame options, used no matter what
	// Stream 0: text frames
	// Stream 1: cursor frames
	args := []string{
		"-y",
		"-r", fmt.Sprint(opts.Framerate),
		"-start_number", fmt.Sprint(opts.StartingFrame),
		"-i", filepath.Join(opts.Input, textFrameFormat),
		"-r", fmt.Sprint(opts.Framerate),
		"-start_number", fmt.Sprint(opts.StartingFrame),
		"-i", filepath.Join(opts.Input, cursorFrameFormat),
	}

	// Add frame image as stream 2 if one is provided
	if opts.ImageFrame != "" {
		args = append(args,
			"-loop", "1",
			"-i", opts.ImageFrame,
		)
	}

	// Build filter code
	var filterArgs strings.Builder
	var prevStageName string

	// Compute dimensions.
	// We do this separately because these
	// values depend on settings.
	var termWidth int
	var termHeight int
	if opts.ImageFrame != "" {
		termWidth = opts.Width - (opts.Padding * 2) - (opts.FrameSize * 2)
		termHeight = opts.Height - (opts.Padding * 2) - (opts.FrameSize * 2)
	} else {
		termWidth = opts.Width - (opts.Padding * 2)
		termHeight = opts.Height - (opts.Padding * 2)
	}

	// The following is used by ALL formats:
	filterArgs.WriteString(
		fmt.Sprintf(`
		[0][1]overlay[merged];
		[merged]scale=%d:%d:force_original_aspect_ratio=1[scaled];
		[scaled]fps=%d,setpts=PTS/%f[speed];
		[speed]pad=%d:%d:(ow-iw)/2:(oh-ih)/2:%s[padded];
		[padded]fillborders=left=%d:right=%d:top=%d:bottom=%d:mode=fixed:color=%s[bordered]
		`,
			termWidth,
			termHeight,

			opts.Framerate,
			opts.PlaybackSpeed,

			termWidth+(opts.Padding+opts.Padding),
			termHeight+(opts.Padding+opts.Padding),
			opts.BackgroundColor,

			opts.Padding,
			opts.Padding,
			opts.Padding,
			opts.Padding,
			opts.BackgroundColor,
		),
	)
	prevStageName = "bordered"

	if opts.ImageFrame != "" {
		// Overlay terminal on background

		// ffmpeg will complain if the final filter ends with a semicolon,
		// so we add one right BEFORE adding additional options.
		filterArgs.WriteString(";")
		filterArgs.WriteString(
			fmt.Sprintf(`
			[2]scale=%d:%d[bg];
			[bg][%s]overlay=(W-w)/2:(H-h)/2:shortest=1[withbg]
			`,
				opts.Width,
				opts.Height,
				prevStageName,
			),
		)
		prevStageName = "withbg"
	}

	// Format-specific options

	if filepath.Ext(targetFile) == ".gif" {
		filterArgs.WriteString(";")
		filterArgs.WriteString(
			fmt.Sprintf(`
			[%s]split[p_a][p_b];
			[p_a]palettegen=max_colors=256[plt];
			[p_b][plt]paletteuse[palette]`,
				prevStageName,
			),
		)
		prevStageName = "palette"
	} else if filepath.Ext(targetFile) == ".webm" {
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-an",
			"-crf", "30",
			"-b:v", "0",
		)
	} else if filepath.Ext(targetFile) == ".mp4" {
		args = append(args,
			"-vcodec", "libx264",
			"-pix_fmt", "yuv420p",
			"-an",
			"-crf", "20",
		)
	}

	args = append(args,
		"-filter_complex", filterArgs.String(),
		"-map", "["+prevStageName+"]",
		targetFile,
	)

	return args
}

// MakeGIF takes a list of images (as frames) and converts them to a GIF.
func MakeGIF(opts VideoOptions) *exec.Cmd {
	targetFile := opts.Output.GIF

	if opts.Output.GIF == "" && opts.Output.WebM == "" && opts.Output.MP4 == "" {
		targetFile = "out.gif"
	} else if opts.Output.GIF == "" {
		return nil
	}

	log.Println(GrayStyle.Render("Creating " + opts.Output.GIF + "..."))
	ensureDir(opts.Output.GIF)

	//nolint:gosec
	return exec.Command(
		"ffmpeg",
		buildFFopts(opts, targetFile)...,
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
