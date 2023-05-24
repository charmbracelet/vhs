package main

import (
	"fmt"
	"os"
	"strings"

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

type StyleOptions struct {
	Width           int
	Height          int
	Padding         int
	BackgroundColor string
	MarginFill      string
	Margin          int
	WindowBar       string
	WindowBarSize   int
	WindowBarColor  string
	BorderRadius    int
}

func DefaultStyleOptions() *StyleOptions {
	return &StyleOptions{
		Width:           defaultWidth,
		Height:          defaultHeight,
		Padding:         defaultPadding,
		MarginFill:      DefaultTheme.Background,
		Margin:          0,
		WindowBar:       "",
		WindowBarSize:   defaultWindowBarSize,
		WindowBarColor:  DefaultTheme.Background,
		BorderRadius:    0,
		BackgroundColor: DefaultTheme.Background,
	}
}

func buildFFMarginFilter(opts *StyleOptions) []string {
	var filter []string

	if opts.MarginFill != "" {
		if marginFillIsColor(opts.MarginFill) {
			// Create plain color stream
			filter = append(filter,
				"-f", "lavfi",
				"-i",
				fmt.Sprintf(
					"color=%s:s=%dx%d",
					opts.MarginFill,
					opts.Width,
					opts.Height,
				),
			)
		} else {
			// Check for existence first.
			_, err := os.Stat(opts.MarginFill)
			if err != nil {
				fmt.Println(ErrorStyle.Render("Unable to read margin file: "), opts.MarginFill)
			}

			// Add image stream
			filter = append(filter,
				"-loop", "1",
				"-i", opts.MarginFill,
			)
		}
	}

	return filter
}

func buildFFBarFilter(opts *StyleOptions, termWidth, termHeight int, barPath string) []string {
	var filter []string

	if opts.WindowBar != "" {
		MakeWindowBar(termWidth, termHeight, *opts, barPath)

		filter = append(filter,
			"-i", barPath,
		)
	}

	return filter
}

func buildFFCornerMarkFilter(opts *StyleOptions, maskPath string, termWidth, termHeight int) []string {
	var filter []string

	if opts.BorderRadius != 0 {
		if opts.WindowBar != "" {
			MakeBorderRadiusMask(termWidth, termHeight+opts.WindowBarSize, opts.BorderRadius, maskPath)
		} else {
			MakeBorderRadiusMask(termWidth, termHeight, opts.BorderRadius, maskPath)
		}

		filter = append(filter,
			"-i", maskPath,
		)
	}

	return filter
}

func addWindowBarFilterCode(filterCode *strings.Builder, opts *StyleOptions, barStream int, prevStageName string) (*strings.Builder, string) {
	if opts.WindowBar != "" {
		// if filterCode.Len() > 0 {
		// filterCode.WriteString(";")
		// }

		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
			[%d]loop=-1[loopbar];
			[loopbar][%s]overlay=0:%d[withbar]
			`,
				barStream,
				prevStageName,
				opts.WindowBarSize,
			),
		)

		prevStageName = "withbar"
	}

	return filterCode, prevStageName
}

func addBorderRadiusFilterCode(filterCode *strings.Builder, opts *StyleOptions, streamId int, prevStageName string) (*strings.Builder, string) {
	if opts.BorderRadius != 0 {
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
				[%d]loop=-1[loopmask];
				[%s][loopmask]alphamerge[rounded]
				`,
				streamId,
				prevStageName,
			),
		)
		prevStageName = "rounded"
	}

	return filterCode, prevStageName
}

func addMarginFillFilterCode(filterCode *strings.Builder, opts *StyleOptions, streamId int, prevStageName string) (*strings.Builder, string) {
	// Overlay terminal on margin
	if opts.MarginFill != "" {
		// ffmpeg will complain if the final filter ends with a semicolon,
		// so we add one BEFORE we start adding filters.
		filterCode.WriteString(";")
		filterCode.WriteString(
			fmt.Sprintf(`
			[%d]scale=%d:%d[bg];
			[bg][%s]overlay=(W-w)/2:(H-h)/2:shortest=1[withbg]
			`,
				streamId,
				opts.Width,
				opts.Height,
				prevStageName,
			),
		)
		prevStageName = "withbg"
	}
	return filterCode, prevStageName
}
