package tty

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Options is the set of options to pass to `ttyd`.
type Options struct {
	Port       int
	FontFamily string
	FontSize   int
	LineHeight float64
	Debug      bool
	Theme      Theme
}

// DefaultOptions are the default options for the `tty`.
func DefaultOptions() Options {
	return Options{
		FontFamily: "SF Mono",
		FontSize:   22,
		LineHeight: 1.2,
		Theme:      DefaultTheme,
	}
}

// Start starts the ttyd process on the given port and options.
func Start(opts Options) *exec.Cmd {
	theme, _ := json.Marshal(opts.Theme)

	cmd := exec.Command(
		"ttyd", fmt.Sprintf("--port=%d", opts.Port),
		"-t", fmt.Sprintf("fontFamily='%s'", opts.FontFamily),
		"-t", fmt.Sprintf("fontSize=%d", opts.FontSize),
		"-t", fmt.Sprintf("lineHeight=%f", opts.LineHeight),
		"-t", fmt.Sprintf("theme=%s", theme),
		"-t", "customGlyphs=true",
		"zsh", "-d", "-f",
	)
	if opts.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
