package ttyd

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Options is the set of options to pass to `ttyd`.
type Options struct {
	Port       int
	FontFamily string
	FontSize   int
	LineHeight float64
}

// DefaultOptions are the default options for the `tty`.
func DefaultOptions() Options {
	return Options{
		FontFamily: "SF Mono",
		FontSize:   22,
		LineHeight: 1.2,
		Port:       7681,
	}
}

// Start starts the ttyd process on the given port and options.
func Start(opts Options) *exec.Cmd {
	theme, _ := json.Marshal(DefaultTheme)

	return exec.Command(
		"ttyd", fmt.Sprintf("--port=%d", opts.Port),
		"-t", fmt.Sprintf("fontFamily='%s'", opts.FontFamily),
		"-t", fmt.Sprintf("fontSize=%d", opts.FontSize),
		"-t", fmt.Sprintf("lineHeight=%f", opts.LineHeight),
		"-t", fmt.Sprintf("theme=%s", string(theme)),
		"-t", "customGlyphs=true",
		"zsh",
	)
}
