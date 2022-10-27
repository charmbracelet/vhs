package main

import (
	"fmt"
	"strings"
)

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

// Underline returns a string of ^ characters which helps underline the problematic token
// in a ParserError.
func Underline(n int) string {
	return ErrorStyle.Render(strings.Repeat("^", n))
}

// LineNumber returns a formatted version of the given line number.
func LineNumber(line int) string {
	return LineNumberStyle.Render(fmt.Sprintf(" %2d │ ", line))
}
