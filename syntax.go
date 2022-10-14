package vhs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/vhs/style"
)

func (c Command) Highlight(faint bool) string {
	var (
		optionsStyle = style.Time
		argsStyle    = style.Number
	)

	if faint {
		if c.Options != "" {
			return style.Faint.Render(fmt.Sprintf("%s %s %s", c.Type, c.Options, c.Args))
		} else {
			return style.Faint.Render(fmt.Sprintf("%s %s", c.Type, c.Args))
		}
	}

	switch c.Type {
	case SET:
		optionsStyle = style.Keyword
		if isNumber(c.Args) {
			argsStyle = style.Number
		} else if isTime(c.Args) {
			argsStyle = style.Time
		} else {
			argsStyle = style.String
		}
	case OUTPUT:
		optionsStyle = style.None
		argsStyle = style.String
	case CTRL:
		argsStyle = style.Command
	case SLEEP:
		argsStyle = style.Time
	case TYPE:
		optionsStyle = style.Time
		argsStyle = style.String
	case HIDE, SHOW:
		return style.Faint.Render(c.Type.String())
	}

	var s strings.Builder
	s.WriteString(style.Command.Render(c.Type.String()) + " ")
	if c.Options != "" {
		s.WriteString(optionsStyle.Render(c.Options) + " ")
	}
	s.WriteString(argsStyle.Render(c.Args))
	return s.String()
}

var numberRegex = regexp.MustCompile("^[0-9]+$")

func isNumber(s string) bool {
	return numberRegex.MatchString(s)
}

var timeRegex = regexp.MustCompile("^[0-9]+m?s$")

func isTime(s string) bool {
	return timeRegex.MatchString(s)
}
