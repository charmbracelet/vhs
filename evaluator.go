package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// EvaluatorOption is a function that can be used to modify the VHS instance.
type EvaluatorOption func(*VHS)

// Evaluate takes as input a tape string, an output writer, and an output file
// and evaluates all the commands within the tape string and produces a GIF.
func Evaluate(tape string, out io.Writer, opts ...EvaluatorOption) error {
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
		if cmd.Type == SET || cmd.Type == OUTPUT {
			fmt.Fprintln(out, cmd.Highlight(false))
			cmd.Execute(&v)
		} else {
			offset = i
			break
		}
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
	v.Record()

	defer v.Cleanup()

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
		if isSetting {
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

	return nil
}
