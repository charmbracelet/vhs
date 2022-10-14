package style

import "github.com/charmbracelet/lipgloss"

var Command = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
var None = lipgloss.NewStyle()
var Faint = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "242", Dark: "238"})
var Keyword = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
var Number = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
var String = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
var Time = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

var Red = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
var Gray = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

var ErrorFile = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("8")).
	Foreground(lipgloss.Color("1")).
	Padding(0, 1).
	Width(80)
