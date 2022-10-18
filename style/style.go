package style

import "github.com/charmbracelet/lipgloss"

const Background = "#171717"
const Foreground = "#dddddd"
const Black = "#282a2e"
const BrightBlack = "#4d4d4d"
const Red = "#c73b1d"
const BrightRed = "#e82100"
const Green = "#00a800"
const BrightGreen = "#00db00"
const Yellow = "#acaf15"
const BrightYellow = "#e5e900"
const Blue = "#3854FC"
const BrightBlue = "#566BF9"
const Magenta = "#d533ce"
const BrightMagenta = "#e83ae9"
const Cyan = "#2cbac9"
const BrightCyan = "#00e6e7"
const White = "#bfbfbf"
const BrightWhite = "#e6e6e6"

const Indigo = "#5B56E0"

var Command = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
var None = lipgloss.NewStyle()
var Faint = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "242", Dark: "238"})
var Keyword = lipgloss.NewStyle()
var Number = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
var String = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
var Time = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
var LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
var Error = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

var File = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

var ErrorFile = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("8")).
	Foreground(lipgloss.Color("1")).
	Padding(0, 1).
	Width(80)
