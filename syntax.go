package vhs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	commandStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	faintStyle   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "242", Dark: "238"})
	numberStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	stringStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	timeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

func (c Command) Highlight(faint bool) string {
	var (
		optionsStyle = timeStyle
		argsStyle    = numberStyle
	)

	if faint {
		if c.Options != "" {
			return faintStyle.Render(fmt.Sprintf("%s %s %s", c.Type, c.Options, c.Args))
		} else {
			return faintStyle.Render(fmt.Sprintf("%s %s", c.Type, c.Args))
		}
	}

	switch c.Type {
	case SET:
		optionsStyle = keywordStyle
		if isNumber(c.Args) {
			argsStyle = numberStyle
		} else {
			argsStyle = stringStyle
		}
	case OUTPUT:
		argsStyle = stringStyle
	case CTRL:
		argsStyle = commandStyle
	case SLEEP:
		argsStyle = timeStyle
	case TYPE:
		optionsStyle = timeStyle
		argsStyle = stringStyle
	case HIDE, SHOW:
		return faintStyle.Render(c.Type.String())
	}

	var s strings.Builder
	s.WriteString(commandStyle.Render(c.Type.String()) + " ")
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
