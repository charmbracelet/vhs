package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// EvaluatorOption is a function that can be used to modify the VHS instance.
type EvaluatorOption func(*VHS)

// Evaluate takes as input a tape string, an output writer, and an output file
// and evaluates all the commands within the tape string and produces a GIF.
func Evaluate(tape string, out io.Writer, opts ...EvaluatorOption) error {
	tape, err := ExecuteTemplate(tape)
	if err != nil {
		return err
	}

	l := NewLexer(tape)
	p := NewParser(l)

	cmds := p.Parse()
	errs := p.Errors()
	if len(errs) != 0 || len(cmds) == 0 {
		lines := strings.Split(tape, "\n")
		for _, err := range errs {
			fmt.Fprint(out, LineNumber(err.Token.Line))
			fmt.Fprintln(out, lines[err.Token.Line-1])
			fmt.Fprint(out, strings.Repeat(" ", err.Token.Column+ErrorColumnOffset))
			fmt.Fprintln(out, Underline(len(err.Token.Literal)), err.Msg)
			fmt.Fprintln(out)
		}
		return errors.New("parse error")
	}

	v := New()
	defer func() { _ = v.close() }()

	// Run Output and Set commands as they only modify options on the VHS instance.
	var offset int
	for i, cmd := range cmds {
		if cmd.Type == SET || cmd.Type == OUTPUT || cmd.Type == REQUIRE {
			fmt.Fprintln(out, cmd.Highlight(false))
			cmd.Execute(&v)
		} else {
			offset = i
			break
		}
	}

	video := v.Options.Video
	if video.Height < 2*video.Padding || video.Width < 2*video.Padding {
		v.Errors = append(v.Errors, fmt.Errorf("height and width must be greater than %d", 2*video.Padding))
	}

	if len(v.Errors) > 0 {
		for _, err := range v.Errors {
			fmt.Fprintln(out, ErrorStyle.Render(err.Error()))
		}
		os.Exit(1)
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
	ctx, cancel := context.WithCancel(context.Background())
	ch := v.Record(ctx)

	// Log errors from the recording process.
	go func() {
		for err := range ch {
			log.Print(err.Error())
		}
	}()

	for _, cmd := range cmds[offset:] {
		// When changing the FontFamily, FontSize, LineHeight, Padding
		// The xterm.js canvas changes dimensions and causes FFMPEG to not work
		// correctly (specifically) with palettegen.
		// It will be possible to change settings on the fly in the future, but it is currently not
		// as it does not result in a proper render of the GIF as the frame sequence
		// will change dimensions. This is fixable.
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

	// Stop recording frames.
	cancel()
	// Read from channel to ensure recorder is done.
	<-ch

	v.Cleanup()
	return nil
}
