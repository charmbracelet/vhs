package main

import (
	"fmt"
	"io"
	"strings"
)

// InvalidSyntaxError is returned when the parser encounters one or more errors.
type InvalidSyntaxError struct {
	Errors []ParserError
}

func (e InvalidSyntaxError) Error() string {
	return fmt.Sprintf("parser: %d error(s)", len(e.Errors))
}

// ParserError represents an error with parsing a tape file.
// It tracks the token causing the error and a human readable error message.
type ParserError struct {
	Token Token
	Msg   string
}

// NewError returns a new ParserError with the given token and message.
func NewError(token Token, msg string) ParserError {
	return ParserError{
		Token: token,
		Msg:   msg,
	}
}

// ErrorColumnOffset is the number of columns that an error should be printed
// to the left to account for the line number.
const ErrorColumnOffset = 5

// String returns a human readable error message printing the token line number
// and message.
func (e ParserError) String() string {
	return fmt.Sprintf("%2d:%-2d │ %s", e.Token.Line, e.Token.Column, e.Msg)
}

func (e ParserError) Error() string {
	return e.String()
}

// Underline returns a string of ^ characters which helps underline the problematic token
// in a ParserError.
func Underline(n int) string {
	return ErrorStyle.Render(strings.Repeat("^", n))
}

// LineNumber returns a formatted version of the given line number.
func LineNumber(line int) string {
	return LineNumberStyle.Render(fmt.Sprintf(" %2d │ ", line))
}

func printParserError(out io.Writer, tape string, err ParserError) {
	lines := strings.Split(tape, "\n")

	fmt.Fprint(out, LineNumber(err.Token.Line))
	fmt.Fprintln(out, lines[err.Token.Line-1])
	fmt.Fprint(out, strings.Repeat(" ", err.Token.Column+ErrorColumnOffset))
	fmt.Fprintln(out, Underline(len(err.Token.Literal)), err.Msg)
	fmt.Fprintln(out)
}

func printErrors(out io.Writer, tape string, errs []error) {
	for _, err := range errs {
		switch err := err.(type) {
		case InvalidSyntaxError:
			for _, v := range err.Errors {
				printParserError(out, tape, v)
			}
			fmt.Fprintln(out, ErrorStyle.Render(err.Error()))

		default:
			fmt.Fprintln(out, ErrorStyle.Render(err.Error()))
		}
	}
}
