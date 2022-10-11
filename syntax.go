package vhs

import (
	"fmt"
	"regexp"

	"github.com/muesli/termenv"
)

var (
	commandColor = termenv.ANSI.Color("12")
	numberColor  = termenv.ANSI.Color("9")
	settingColor = termenv.ANSI.Color("14")
	stringColor  = termenv.ANSI.Color("10")
	timeColor    = termenv.ANSI.Color("3")
	pathColor    = termenv.ANSI.Color("2")
)

func (c Command) Highlight() string {
	command := termenv.String(c.Type.String()).Foreground(commandColor)

	var options = termenv.String(c.Options)
	var args = termenv.String(c.Args)

	switch c.Type {
	case SET:
		options = options.Foreground(settingColor)
		if isNumber(c.Args) {
			args = args.Foreground(numberColor)
		} else {
			args = args.Foreground(stringColor)
		}
	case OUTPUT:
		args = args.Foreground(pathColor)
	case CTRL:
		args = args.Foreground(commandColor)
	case SLEEP:
		args = args.Foreground(timeColor)
	case TYPE:
		args = args.Foreground(stringColor)
		options = options.Foreground(timeColor)
	default:
		options = options.Foreground(timeColor)
		args = args.Foreground(numberColor)
	}

	if c.Options != "" {
		return fmt.Sprintf("%s %s %s", command, options, args)
	}
	return fmt.Sprintf("%s %s", command, args)
}

var numberRegex = regexp.MustCompile("^[0-9]+$")

func isNumber(s string) bool {
	return numberRegex.MatchString(s)
}
