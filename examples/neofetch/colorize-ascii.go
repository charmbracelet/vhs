package neofetch

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

//go:embed vhs.ascii
var art string

func main() {
	lipgloss.SetColorProfile(termenv.TrueColor)

	var (
		b      strings.Builder
		lines  = strings.Split(art, "\n")
		colors = []string{"#AD92FF", "#9D7CFF", "#906BFF", "#8056FF"}
		step   = len(lines) / len(colors)
	)

	for i, l := range lines {
		n := clamp(0, len(colors)-1, i/step)
		b.WriteString(colorize(colors[n], l))
		b.WriteRune('\n')
	}

	fmt.Print(b.String())
}

func colorize(c, s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render(s)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
