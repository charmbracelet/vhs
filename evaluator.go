package vhs

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func Evaluate(tape string, w io.Writer) error {
	v := New()
	defer v.Cleanup()

	l := NewLexer(tape)
	p := NewParser(l)

	cmds := p.Parse()
	errs := p.Errors()
	if len(errs) != 0 {
		lines := strings.Split(tape, "\n")
		for _, err := range errs {
			fmt.Fprint(w, LineNumber(err.Token.Line))
			fmt.Fprintln(w, lines[err.Token.Line-1])
			fmt.Fprint(w, strings.Repeat(" ", err.Token.Column+5))
			fmt.Fprintln(w, Underline(len(err.Token.Literal)), err.Msg)
			fmt.Fprintln(w)
		}
		return errors.New("Parse error")
	}

	var offset int

	for i, cmd := range cmds {
		if cmd.Type == Set {
			fmt.Fprintf(w, "Setting %s to %s\n", cmd.Options, cmd.Args)
			cmd.Execute(&v)
		} else {
			offset = i
			break
		}
	}

	v.Start()

	for _, cmd := range cmds[offset:] {
		fmt.Fprintln(w, cmd)
		cmd.Execute(&v)
	}

	return nil
}
