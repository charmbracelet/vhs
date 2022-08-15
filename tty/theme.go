package tty

// Theme is a terminal theme for xterm
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

// DefaultTheme is the default theme to use for recording demos and
// screenshots. Taken from https://github.com/meowgorithm/dotfiles.
var DefaultTheme = Theme{
	Background:    "#171717",
	Black:         "#000000",
	Blue:          "#3854FC",
	BrightBlack:   "#4d4d4d",
	BrightBlue:    "#566BF9",
	BrightCyan:    "#00e6e7",
	BrightGreen:   "#00db00",
	BrightMagenta: "#e83ae9",
	BrightRed:     "#e82100",
	BrightWhite:   "#e6e6e6",
	BrightYellow:  "#e5e900",
	Cyan:          "#2cbac9",
	Foreground:    "#dddddd",
	Green:         "#00a800",
	Magenta:       "#d533ce",
	Red:           "#c73b1d",
	White:         "#bfbfbf",
	Yellow:        "#acaf15",
}
