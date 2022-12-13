package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
)

// EvaluatorOption is a function that can be used to modify the VHS instance.
type EvaluatorOption func(*VHS)

// Evaluate takes as input a tape string, an output writer, and an output file
// and evaluates all the commands within the tape string and produces a GIF.
func Evaluate(ctx context.Context, tape string, out io.Writer, opts ...EvaluatorOption) []error {
	l := NewLexer(tape)
	p := NewParser(l)

	cmds := p.Parse()
	errs := p.Errors()
	if len(errs) != 0 || len(cmds) == 0 {
		return []error{InvalidSyntaxError{errs}}
	}

	v := New()
	for _, cmd := range cmds {
		if cmd.Type == SET && cmd.Options == "Shell" {
			cmd.Execute(&v)
		}
	}

	// Start things up
	if err := v.Start(); err != nil {
		return []error{err}
	}
	defer func() { _ = v.close() }()

	// Run Output and Set commands as they only modify options on the VHS instance.
	var offset int
	for i, cmd := range cmds {
		if cmd.Type == SET || cmd.Type == OUTPUT || cmd.Type == REQUIRE {
			fmt.Fprintln(out, cmd.Highlight(false))
			if cmd.Options != "Shell" {
				cmd.Execute(&v)
			}
		} else {
			offset = i
			break
		}
	}

	video := v.Options.Video
	minDimension := (2 * video.Padding) + (2 * video.MarginSize)
	if video.Height < minDimension || video.Width < minDimension {
		v.Errors = append(v.Errors, fmt.Errorf("width and height must be greater than %d to account for padding and frames", minDimension))
	}

	// Check if margin color is a valid hex string
	if video.MarginIsColor {
		_, err := strconv.ParseUint(
			video.MarginFill[1:],
			16,
			64,
		)

		if err != nil || len(video.MarginFill) != 7 {
			v.Errors = append(
				v.Errors,
				fmt.Errorf(
					"\"%s\" is not a valid color",
					video.MarginFill,
				),
			)
		}
	}

	if len(v.Errors) > 0 {
		return v.Errors
	}

	// Setup the terminal session so we can start executing commands.
	v.Setup()

	// If the first command (after Settings and Outputs) is a Hide command, we can
	// begin executing the commands before we start recording to avoid capturing
	// any unwanted frames.
	if cmds[offset].Type == HIDE {
		for i, cmd := range cmds[offset:] {
			if cmd.Type == SHOW {
				offset += i
				break
			}
			fmt.Fprintln(out, cmd.Highlight(true))
			cmd.Execute(&v)
		}
	}

	// Begin recording frames as we are now in a recording state.
	ctx, cancel := context.WithCancel(ctx)
	ch := v.Record(ctx)

	// Clean up temporary files at the end.
	defer func() { _ = v.Cleanup() }()

	teardown := func() {
		// Stop recording frames.
		cancel()
		// Read from channel to ensure recorder is done.
		<-ch
	}

	// Log errors from the recording process.
	go func() {
		for err := range ch {
			log.Print(err.Error())
		}
	}()

	for _, cmd := range cmds[offset:] {
		if ctx.Err() != nil {
			teardown()
			return []error{ctx.Err()}
		}

		// When changing the FontFamily, FontSize, LineHeight, Padding
		// The xterm.js canvas changes dimensions and causes FFMPEG to not work
		// correctly (specifically) with palettegen.
		// It will be possible to change settings on the fly in the future, but
		// it is currently not as it does not result in a proper render of the
		// GIF as the frame sequence will change dimensions. This is fixable.
		//
		// We should remove if isSetting statement.
		isSetting := cmd.Type == SET && cmd.Options != "TypingSpeed"
		if isSetting || cmd.Type == REQUIRE {
			fmt.Fprintln(out, cmd.Highlight(true))
			continue
		}
		fmt.Fprintln(out, cmd.Highlight(!v.recording || cmd.Type == SHOW || cmd.Type == HIDE || isSetting))
		cmd.Execute(&v)
	}

	// If running as an SSH server, the output file is a temporary file
	// to use for the output.
	//
	// We need to set the GIF file path before it is created but after all of
	// the settings and commands are executed. This is done in `serve.go`.
	//
	// Since the GIF creation is deferred, setting the output file here will
	// achieve what we want.
	for _, opt := range opts {
		opt(&v)
	}

	teardown()
	if err := v.Render(); err != nil {
		return []error{err}
	}
	return nil
}
