package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/vhs/parser"
)

// InvalidSyntaxError is returned when the parser encounters one or more errors.
type InvalidSyntaxError struct {
	Errors []parser.Error
}

func (e InvalidSyntaxError) Error() string {
	return fmt.Sprintf("parser: %d error(s)", len(e.Errors))
}

// ErrorColumnOffset is the number of columns that an error should be printed
// to the left to account for the line number.
const ErrorColumnOffset = 5

// Underline returns a string of ^ characters which helps underline the problematic token
// in a parser.Error.
func Underline(n int) string {
	return ErrorStyle.Render(strings.Repeat("^", n))
}

// LineNumber returns a formatted version of the given line number.
func LineNumber(line int) string {
	return LineNumberStyle.Render(fmt.Sprintf(" %2d │ ", line))
}

func printError(out io.Writer, tape string, err parser.Error) {
	lines := strings.Split(tape, "\n")

	_, _ = fmt.Fprint(out, LineNumber(err.Token.Line))
	_, _ = fmt.Fprintln(out, lines[err.Token.Line-1])
	_, _ = fmt.Fprint(out, strings.Repeat(" ", err.Token.Column+ErrorColumnOffset))
	_, _ = fmt.Fprintln(out, Underline(len(err.Token.Literal)), err.Msg)
	_, _ = fmt.Fprintln(out)
}

func printErrors(out io.Writer, tape string, errs []error) {
	for _, err := range errs {
		switch err := err.(type) {
		case InvalidSyntaxError:
			for _, v := range err.Errors {
				printError(out, tape, v)
			}
			_, _ = fmt.Fprintln(out, ErrorStyle.Render(err.Error()))

		default:
			_, _ = fmt.Fprintln(out, ErrorStyle.Render(err.Error()))
		}
	}
}
