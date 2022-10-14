package vhs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

var redStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
var grayStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

func Underline(n int) string {
	return redStyle.Render(strings.Repeat("^", n))
}

func LineNumber(line int) string {
	return grayStyle.Render(fmt.Sprintf(" %2d │ ", line))
}
