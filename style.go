package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme colors.
const (
	Background    = "#171717"
	Foreground    = "#dddddd"
	Black         = "#282a2e" // ansi 0
	BrightBlack   = "#4d4d4d" // ansi 8
	Red           = "#D74E6F" // ansi 1
	BrightRed     = "#FE5F86" // ansi 9
	Green         = "#31BB71" // ansi 2
	BrightGreen   = "#00D787" // ansi 10
	Yellow        = "#D3E561" // ansi 3
	BrightYellow  = "#EBFF71" // ansi 11
	Blue          = "#8056FF" // ansi 4
	BrightBlue    = "#9B79FF" // ansi 12
	Magenta       = "#ED61D7" // ansi 5
	BrightMagenta = "#FF7AEA" // ansi 13
	Cyan          = "#04D7D7" // ansi 6
	BrightCyan    = "#00FEFE" // ansi 14
	White         = "#bfbfbf" // ansi 7
	BrightWhite   = "#e6e6e6" // ansi 15
	Indigo        = "#5B56E0"
)

const (
	defaultColumns       = 80
	defaultHeight        = 600
	defaultMaxColors     = 256
	defaultPadding       = 60
	defaultWindowBarSize = 30
	defaultPlaybackSpeed = 1.0
	defaultWidth         = 1200
)

// Styles for syntax highlighting
var (
	CommandStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	FaintStyle      = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "242", Dark: "238"})
	NoneStyle       = lipgloss.NewStyle()
	KeywordStyle    = lipgloss.NewStyle()
	URLStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	NumberStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	StringStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	TimeStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	LineNumberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	ErrorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	GrayStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	ErrorFileStyle  = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8")).
			Foreground(lipgloss.Color("1")).
			Padding(0, 1).
			Width(defaultColumns)
)

// StyleOptions represents the ui options for video and screenshots.
type StyleOptions struct {
	Width            int
	Height           int
	Padding          int
	BackgroundColor  string
	MarginFill       string
	Margin           int
	WindowBar        string
	WindowTitle      string
	WindowTitleColor string
	WindowBarSize    int
	WindowBarColor   string
	BorderRadius     int
}

// DefaultStyleOptions returns default Style config.
func DefaultStyleOptions() *StyleOptions {
	return &StyleOptions{
		Width:            defaultWidth,
		Height:           defaultHeight,
		Padding:          defaultPadding,
		MarginFill:       DefaultTheme.Background,
		Margin:           0,
		WindowBar:        "",
		WindowTitle:      "",
		WindowTitleColor: DefaultTheme.Foreground,
		WindowBarSize:    defaultWindowBarSize,
		WindowBarColor:   DefaultTheme.Background,
		BorderRadius:     0,
		BackgroundColor:  DefaultTheme.Background,
	}
}
