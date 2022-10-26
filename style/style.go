package style

import "github.com/charmbracelet/lipgloss"

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

const defaultWidth = 80

// Command is the default style for command keywords
var Command = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

// None is an empty style
var None = lipgloss.NewStyle()

// Faint is the default style for a faint command
var Faint = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "242", Dark: "238"})

// Keyword is the style for a keyword
var Keyword = lipgloss.NewStyle()

// Number is the style for a number
var Number = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

// String is the style for a string
var String = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

// Time is the style for time
var Time = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

// LineNumber is the style for a line number
var LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

// Error is the style for an error
var Error = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

// File is the style for printing a file name
var File = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

// ErrorFile is the style for printing an file name that is not valid.
var ErrorFile = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("8")).
	Foreground(lipgloss.Color("1")).
	Padding(0, 1).
	Width(defaultWidth)
