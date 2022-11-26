package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Highlight syntax highlights a command for prettier printing.
// It takes an argument whether or not to print the command in a faint style to
// represent hidden commands.
func (c Command) Highlight(faint bool) string {
	var (
		optionsStyle = TimeStyle
		argsStyle    = NumberStyle
	)

	if faint {
		if c.Options != "" {
			return FaintStyle.Render(fmt.Sprintf("%s %s %s", c.Type, c.Options, c.Args))
		}
		return FaintStyle.Render(fmt.Sprintf("%s %s", c.Type, c.Args))
	}

	switch c.Type {
	case MATCH, MATCH_ANY:
		argsStyle = StringStyle
	case SET:
		optionsStyle = KeywordStyle
		if isNumber(c.Args) {
			argsStyle = NumberStyle
		} else if isTime(c.Args) {
			argsStyle = TimeStyle
		} else {
			argsStyle = StringStyle
		}
	case OUTPUT:
		optionsStyle = NoneStyle
		argsStyle = StringStyle
	case CTRL:
		argsStyle = CommandStyle
	case SLEEP:
		argsStyle = TimeStyle
	case TYPE:
		optionsStyle = TimeStyle
		argsStyle = StringStyle
	case HIDE, SHOW:
		return FaintStyle.Render(c.Type.String())
	}

	var s strings.Builder
	s.WriteString(CommandStyle.Render(c.Type.String()) + " ")
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
