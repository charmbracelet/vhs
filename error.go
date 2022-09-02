package vhs

import (
	"fmt"
	"strings"
)

type ParserError struct {
	Token Token
	Msg   string
}

func NewError(token Token, msg string) ParserError {
	return ParserError{
		Token: token,
		Msg:   msg,
	}
}

func (e ParserError) String() string {
	return fmt.Sprintf("%2d:%-2d │ %s", e.Token.Line, e.Token.Column, e.Msg)
}

func Underline(n int) string {
	return fmt.Sprintf("%s%s%s", "\x1b[31m", strings.Repeat("^", n), "\x1b[0m")
}

func LineNumber(line int) string {
	return fmt.Sprintf("\x1b[90m %2d │ \x1b[0m", line)
}
