package vhs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/vhs/style"
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
	return style.Error.Render(strings.Repeat("^", n))
}

func LineNumber(line int) string {
	return style.LineNumber.Render(fmt.Sprintf(" %2d │ ", line))
}
