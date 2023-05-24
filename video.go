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

// buildFFopts assembles an ffmpeg command from some VideoOptions
func buildFFopts(opts VideoOptions, targetFile string) []string {
	// Variables used for building ffmpeg command
	var filterCode strings.Builder
	var args []string
	var prevStageName string
	streamCounter := 2

	// Compute dimensions of terminal
	termWidth := opts.Style.Width
	termHeight := opts.Style.Height
	if opts.Style.MarginFill != "" {
		termWidth = termWidth - double(opts.Style.Margin)
		termHeight = termHeight - double(opts.Style.Margin)
	}
	if opts.Style.WindowBar != "" {
		termHeight = termHeight - opts.Style.WindowBarSize
	}

	// Input frame options, used no matter what
	// Stream 0: text frames
	// Stream 1: cursor frames
	args = append(args,
		"-y",
		"-r", fmt.Sprint(opts.Framerate),
		"-start_number", fmt.Sprint(opts.StartingFrame),
		"-i", filepath.Join(opts.Input, textFrameFormat),
		"-r", fmt.Sprint(opts.Framerate),
		"-start_number", fmt.Sprint(opts.StartingFrame),
		"-i", filepath.Join(opts.Input, cursorFrameFormat),
	)

	// Add margin stream if one is provided
	var marginStream int
	if opts.Style.MarginFill != "" {
		if marginFillIsColor(opts.Style.MarginFill) {
			// Create plain color stream
			args = append(args,
				"-f", "lavfi",
				"-i",
				fmt.Sprintf(
					"color=%s:s=%dx%d",
					opts.Style.MarginFill,
					opts.Style.Width,
					opts.Style.Height,
				),
			)
			marginStream = streamCounter
			streamCounter++
		} else {
			// Check for existence first.
			_, err := os.Stat(opts.Style.MarginFill)
			if err != nil {
				fmt.Println(ErrorStyle.Render("Unable to read margin file: "), opts.Style.MarginFill)
			}

			// Add image stream
			args = append(args,
				"-loop", "1",
				"-i", opts.Style.MarginFill,
			)
			marginStream = streamCounter
			streamCounter++
		}
	}

	// Create and add a window bar stream if necessary
	var barStream int
	if opts.Style.WindowBar != "" {
		barPath := filepath.Join(opts.Input, "bar.png")
		MakeWindowBar(termWidth, termHeight, *opts.Style, barPath)

		args = append(args,
			"-i", barPath,
		)
		barStream = streamCounter
		streamCounter++
	}

	// Create and add rounded-corner mask if necessary
	var cornerMaskStream int
	if opts.Style.BorderRadius != 0 {
		borderMaskPath := filepath.Join(opts.Input, "mask.png")
		if opts.Style.WindowBar != "" {
			MakeBorderRadiusMask(termWidth, termHeight+opts.Style.WindowBarSize, opts.Style.BorderRadius, borderMaskPath)
		} else {
			MakeBorderRadiusMask(termWidth, termHeight, opts.Style.BorderRadius, borderMaskPath)
		}

		args = append(args,
			"-i", borderMaskPath,
		)
		cornerMaskStream = streamCounter
		//streamCounter++
	}

	// The following filters are always used
	filterCode.WriteString(
		fmt.Sprintf(`
		[0][1]overlay[merged];
		[merged]scale=%d:%d:force_original_aspect_ratio=1[scaled];
		[scaled]fps=%d,setpts=PTS/%f[speed];
		[speed]pad=%d:%d:(ow-iw)/2:(oh-ih)/2:%s[padded];
		[padded]fillborders=left=%d:right=%d:top=%d:bottom=%d:mode=fixed:color=%s[padded]
		`,
			termWidth-double(opts.Style.Padding),
			termHeight-double(opts.Style.Padding),

			opts.Framerate,
			opts.PlaybackSpeed,

			termWidth,
			termHeight,
			opts.Style.BackgroundColor,

			opts.Style.Padding,
			opts.Style.Padding,
			opts.Style.Padding,
			opts.Style.Padding,
			opts.Style.BackgroundColor,
		),
	)
	prevStageName = "padded"

	// Add a bar to the terminal and mask the output.
	// This allows us to round the corners of the terminal.
	if opts.Style.WindowBar != "" {
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
			[%d]loop=-1[loopbar];
			[loopbar][%s]overlay=0:%d[withbar]
			`,
				barStream,
				prevStageName,
				opts.Style.WindowBarSize,
			),
		)
		prevStageName = "withbar"
	}

	if opts.Style.BorderRadius != 0 {
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
	if opts.Style.MarginFill != "" {
		// ffmpeg will complain if the final filter ends with a semicolon,
		// so we add one BEFORE we start adding filters.
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
			[%d]scale=%d:%d[bg];
			[bg][%s]overlay=(W-w)/2:(H-h)/2:shortest=1[withbg]
			`,
				marginStream,
				opts.Style.Width,
				opts.Style.Height,
				prevStageName,
			),
		)
		prevStageName = "withbg"
	}

	// Format-specific options
	if filepath.Ext(targetFile) == gif {
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
			[%s]split[plt_a][plt_b];
			[plt_a]palettegen=max_colors=256[plt];
			[plt_b][plt]paletteuse[palette]`,
				prevStageName,
			),
		)
		prevStageName = "palette"
	} else if filepath.Ext(targetFile) == webm {
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-an",
			"-crf", "30",
			"-b:v", "0",
		)
	} else if filepath.Ext(targetFile) == mp4 {
		args = append(args,
			"-vcodec", "libx264",
			"-pix_fmt", "yuv420p",
			"-an",
			"-crf", "20",
		)
	}

	args = append(args,
		"-filter_complex", filterCode.String(),
		"-map", "["+prevStageName+"]",
		targetFile,
	)

	return args
}

// MakeGIF takes a list of images (as frames) and converts them to a GIF. func MakeGIF(opts VideoOptions) *exec.Cmd {
func MakeGIF(opts VideoOptions) *exec.Cmd {
	targetFile := opts.Output.GIF

	if opts.Output.GIF == "" && opts.Output.WebM == "" && opts.Output.MP4 == "" {
		targetFile = "out.gif"
	} else if opts.Output.GIF == "" {
		return nil
	}

	log.Println(GrayStyle.Render("Creating " + opts.Output.GIF + "..."))
	ensureDir(opts.Output.GIF)

	fmt.Println(buildFFopts(opts, targetFile))

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
