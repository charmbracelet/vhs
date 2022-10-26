package style

import "github.com/charmbracelet/lipgloss"

// Background is the default background color of the theme.
const Background = "#171717"

// Foreground is the default foreground color of the theme.
const Foreground = "#dddddd"

// Black is the default black color of the theme.
const Black = "#282a2e"

// BrightBlack is the default BrightBlack color of the theme.
const BrightBlack = "#4d4d4d"

// Red is the default Red color of the theme.
const Red = "#D74E6F"

// BrightRed is the default BrightRed color of the theme.
const BrightRed = "#FE5F86"

// Green is the default Green color of the theme.
const Green = "#31BB71"

// BrightGreen is the default BrightGreen color of the theme.
const BrightGreen = "#00D787"

// Yellow is the default Yellow color of the theme.
const Yellow = "#D3E561"

// BrightYellow is the default BrightYellow color of the theme.
const BrightYellow = "#EBFF71"

// Blue is the default Blue color of the theme.
const Blue = "#8056FF"

// BrightBlue is the default BrightBlue color of the theme.
const BrightBlue = "#9B79FF"

// Magenta is the default Magenta color of the theme.
const Magenta = "#ED61D7"

// BrightMagenta is the default BrightMagenta color of the theme.
const BrightMagenta = "#FF7AEA"

// Cyan is the default Cyan color of the theme.
const Cyan = "#04D7D7"

// BrightCyan is the default BrightCyan color of the theme.
const BrightCyan = "#00FEFE"

// White is the default White color of the theme.
const White = "#bfbfbf"

// BrightWhite is the default BrightWhite color of the theme.
const BrightWhite = "#e6e6e6"

// Indigo is the default Indigo color of the theme.
const Indigo = "#5B56E0"

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
