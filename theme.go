// theme contains the information about a terminal theme.
// It stores the 16 base colors as well as the background and foreground colors
// of the terminal theme.
//
// It can be changed through the Set command.
//
// Set Theme {"background": "#171717"}
//
package vhs

import (
	"encoding/json"

	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/vhs/style"
)

// Theme is a terminal theme for xterm.js
// It is used for marshalling between the xterm.js readable json format and a
// valid go struct.
type Theme struct {
	Background    string `json:"background"`
	Foreground    string `json:"foreground"`
	Black         string `json:"black"`
	BrightBlack   string `json:"brightBlack"`
	Red           string `json:"red"`
	BrightRed     string `json:"brightRed"`
	Green         string `json:"green"`
	BrightGreen   string `json:"brightGreen"`
	Yellow        string `json:"yellow"`
	BrightYellow  string `json:"brightYellow"`
	Blue          string `json:"blue"`
	BrightBlue    string `json:"brightBlue"`
	Magenta       string `json:"magenta"`
	BrightMagenta string `json:"brightMagenta"`
	Cyan          string `json:"cyan"`
	BrightCyan    string `json:"brightCyan"`
	White         string `json:"white"`
	BrightWhite   string `json:"brightWhite"`
}

func (t Theme) String() string {
	ts, err := json.Marshal(t)
	if err != nil {
		dts, _ := json.Marshal(DefaultTheme)
		return string(dts)
	}
	return string(ts)
}

// DefaultTheme is the default theme to use for recording demos and
// screenshots.
//
// Taken from https://github.com/meowgorithm/dotfiles.
var DefaultTheme = Theme{
	Background:    style.Background,
	Foreground:    style.Foreground,
	Black:         style.Black,
	BrightBlack:   style.BrightBlack,
	Red:           style.Red,
	BrightRed:     style.BrightRed,
	Green:         style.Green,
	BrightGreen:   style.BrightGreen,
	Yellow:        style.Yellow,
	BrightYellow:  style.BrightYellow,
	Blue:          style.Blue,
	BrightBlue:    style.BrightBlue,
	Magenta:       style.Magenta,
	BrightMagenta: style.BrightMagenta,
	Cyan:          style.Cyan,
	BrightCyan:    style.BrightCyan,
	White:         style.White,
	BrightWhite:   style.BrightWhite,
}

// GlamourTheme is the theme for printing out the manual page
// $ vhs man
var GlamourTheme = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:       stringPtr(style.BrightWhite),
			BlockPrefix: "\n",
			BlockSuffix: "\n",
		},
		Margin: uintPtr(2),
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       stringPtr(style.Indigo),
			Bold:        boolPtr(true),
		},
	},
	Item:     ansi.StylePrimitive{Prefix: "· "},
	Emph:     ansi.StylePrimitive{Color: stringPtr(style.BrightBlack)},
	Strong:   ansi.StylePrimitive{Bold: boolPtr(true)},
	Link:     ansi.StylePrimitive{Color: stringPtr(style.BrightGreen), Underline: boolPtr(true)},
	LinkText: ansi.StylePrimitive{Color: stringPtr(style.BrightMagenta)},
	Code:     ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Color: stringPtr(style.BrightMagenta)}},
}

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }
func color(c string) ansi.StylePrimitive {
	return ansi.StylePrimitive{Color: &c}
}
