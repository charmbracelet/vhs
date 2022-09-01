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
	Background:    "#171717",
	Foreground:    "#dddddd",
	Black:         "#000000",
	BrightBlack:   "#4d4d4d",
	Red:           "#c73b1d",
	BrightRed:     "#e82100",
	Green:         "#00a800",
	BrightGreen:   "#00db00",
	Yellow:        "#acaf15",
	BrightYellow:  "#e5e900",
	Blue:          "#3854FC",
	BrightBlue:    "#566BF9",
	Magenta:       "#d533ce",
	BrightMagenta: "#e83ae9",
	Cyan:          "#2cbac9",
	BrightCyan:    "#00e6e7",
	White:         "#bfbfbf",
	BrightWhite:   "#e6e6e6",
}
