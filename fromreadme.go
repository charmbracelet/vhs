package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// fromReadmeOptions holds the configuration for the from-readme command.
type fromReadmeOptions struct {
	section     string
	command     string
	output      string
	pause       time.Duration
	typingSpeed time.Duration
	dryRun      bool
	tapeOut     string
	fontSize    int
	width       int
	height      int
	waitTimeout time.Duration
	waitPattern string
}

// defaultFromReadmeWaitPattern matches the end of common bash/zsh/fish/root prompts.
// VHS's built-in default (">$") only matches fish; this covers "$", "#", ">", "%".
const defaultFromReadmeWaitPattern = `[$#>%] *$`

var fromReadmeOpts fromReadmeOptions

var fromReadmeCmd = &cobra.Command{
	Use:          "from-readme [file]",
	Short:        "Generate a GIF by running shell commands from a Markdown document",
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
	RunE:         runFromReadme,
}

func init() {
	fromReadmeCmd.Flags().StringVar(&fromReadmeOpts.section, "section", "", "only extract commands under this heading (case-insensitive)")
	fromReadmeCmd.Flags().StringVar(&fromReadmeOpts.command, "command", "", "only include lines whose first token matches this CLI name")
	fromReadmeCmd.Flags().StringVarP(&fromReadmeOpts.output, "output", "o", "vhs.gif", "output GIF file path")
	fromReadmeCmd.Flags().DurationVar(&fromReadmeOpts.pause, "pause", 2*time.Second, "sleep duration between commands")
	fromReadmeCmd.Flags().DurationVar(&fromReadmeOpts.typingSpeed, "typing-speed", 75*time.Millisecond, "per-keystroke typing delay")
	fromReadmeCmd.Flags().BoolVar(&fromReadmeOpts.dryRun, "dry-run", false, "print the generated tape without recording")
	fromReadmeCmd.Flags().StringVar(&fromReadmeOpts.tapeOut, "tape-out", "", "also save the generated tape to this file")
	fromReadmeCmd.Flags().IntVar(&fromReadmeOpts.fontSize, "font-size", 15, "terminal font size")
	fromReadmeCmd.Flags().IntVar(&fromReadmeOpts.width, "width", 1600, "terminal width in pixels")
	fromReadmeCmd.Flags().IntVar(&fromReadmeOpts.height, "height", 900, "terminal height in pixels")
	fromReadmeCmd.Flags().DurationVar(&fromReadmeOpts.waitTimeout, "wait-timeout", 2*time.Minute, "max time to wait for shell prompt after each command (increase for long-running commands)")
	fromReadmeCmd.Flags().StringVar(&fromReadmeOpts.waitPattern, "wait-pattern", defaultFromReadmeWaitPattern, "regex to detect shell prompt (used in Wait commands)")
}

func runFromReadme(cmd *cobra.Command, args []string) error {
	var mdFile string
	if len(args) > 0 {
		mdFile = args[0]
	} else {
		var err error
		mdFile, err = findREADME(".")
		if err != nil {
			return err
		}
		log.Println(GrayStyle.Render("File: " + mdFile))
	}

	content, err := os.ReadFile(mdFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", mdFile, err)
	}

	commands, err := extractShellCommands(string(content), fromReadmeOpts.section)
	if err != nil {
		return err
	}

	if fromReadmeOpts.command != "" {
		commands = filterByCommand(commands, fromReadmeOpts.command)
	}

	if len(commands) == 0 {
		return errors.New("no shell commands found")
	}

	tape := buildFromReadmeTape(commands, fromReadmeOpts)

	if fromReadmeOpts.tapeOut != "" {
		if err := os.WriteFile(fromReadmeOpts.tapeOut, []byte(tape), 0o600); err != nil {
			return fmt.Errorf("writing tape file: %w", err)
		}
		log.Println(GrayStyle.Render("Tape: " + fromReadmeOpts.tapeOut))
	}

	if fromReadmeOpts.dryRun {
		fmt.Fprint(cmd.OutOrStdout(), tape)
		return nil
	}

	if err := ensureDependencies(); err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if quietFlag {
		out = os.Stderr
	}

	errs := Evaluate(cmd.Context(), tape, out)
	if len(errs) > 0 {
		printErrors(os.Stderr, tape, errs)
		return errors.New("recording failed")
	}

	return nil
}

// findREADME searches dir for a README file, returning the first one found.
func findREADME(dir string) (string, error) {
	candidates := []string{"README.md", "Readme.md", "readme.md", "README"}
	for _, name := range candidates {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("no README found in %s", dir)
}

// shellLanguages is the set of fenced code block language tags treated as shell.
var shellLanguages = map[string]bool{
	"sh":         true,
	"bash":       true,
	"zsh":        true,
	"fish":       true,
	"shell":      true,
	"powershell": true,
	"pwsh":       true,
	"cmd":        true,
}

// extractShellCommands parses Markdown content and returns shell command lines
// found in fenced code blocks. If section is non-empty, only blocks under the
// matching heading are included (case-insensitive match). The section ends when
// a heading at the same or higher level is encountered.
func extractShellCommands(content, section string) ([]string, error) {
	lines := strings.Split(content, "\n")

	inSection := section == ""
	sectionLevel := 0
	inCodeBlock := false
	isShellBlock := false

	var commands []string

	for _, line := range lines {
		// Detect fenced code block boundaries (``` or ~~~).
		if isFence(line) {
			if !inCodeBlock {
				lang := strings.ToLower(strings.TrimSpace(strings.TrimLeft(line, "`~")))
				inCodeBlock = true
				isShellBlock = shellLanguages[lang] && inSection
			} else {
				inCodeBlock = false
				isShellBlock = false
			}
			continue
		}

		// While inside a code block, collect non-empty, non-comment lines.
		if inCodeBlock {
			if isShellBlock && line != "" && !strings.HasPrefix(line, "#") {
				commands = append(commands, line)
			}
			continue
		}

		// Handle headings for section filtering.
		if strings.HasPrefix(line, "#") && section != "" {
			level := headingLevel(line)
			text := strings.ToLower(strings.TrimSpace(strings.TrimLeft(line, "# ")))

			switch {
			case strings.Contains(text, strings.ToLower(section)):
				inSection = true
				sectionLevel = level
			case inSection && level <= sectionLevel:
				inSection = false
			}
		}
	}

	return commands, nil
}

// isFence returns true if line is a fenced code block delimiter (``` or ~~~).
func isFence(line string) bool {
	return strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~")
}

// headingLevel returns the ATX heading level (number of leading '#') of line.
func headingLevel(line string) int {
	level := 0
	for _, ch := range line {
		if ch == '#' {
			level++
		} else {
			break
		}
	}
	return level
}

// prefixTokens are command-line wrappers that precede the actual tool name.
var prefixTokens = map[string]bool{
	"sudo": true,
	"npx":  true,
	"env":  true,
	"time": true,
}

// filterByCommand returns only command lines whose effective first token
// (after stripping prefix wrappers like sudo/npx) matches command.
// If command is empty, all lines are returned unchanged.
func filterByCommand(commands []string, command string) []string {
	if command == "" {
		return commands
	}
	var out []string
	for _, line := range commands {
		tokens := strings.Fields(line)
		for len(tokens) > 0 && prefixTokens[tokens[0]] {
			tokens = tokens[1:]
		}
		if len(tokens) == 0 {
			continue
		}
		first := filepath.Base(tokens[0])
		if first == command {
			out = append(out, line)
		}
	}
	return out
}

// tapeDuration formats a time.Duration as the smallest clean VHS tape unit
// (e.g. 2m, 30s, 500ms) without trailing zero components that the parser rejects.
func tapeDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}
	if d%time.Minute == 0 {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d%time.Second == 0 {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}

// buildFromReadmeTape generates a VHS tape string from the given shell commands.
func buildFromReadmeTape(commands []string, opts fromReadmeOptions) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Output %s\n\n", opts.output)
	fmt.Fprintf(&b, "Set FontSize %d\n", opts.fontSize)
	fmt.Fprintf(&b, "Set Width %d\n", opts.width)
	fmt.Fprintf(&b, "Set Height %d\n", opts.height)
	fmt.Fprintf(&b, "Set TypingSpeed %s\n", tapeDuration(opts.typingSpeed))
	fmt.Fprintf(&b, "Set WaitTimeout %s\n", tapeDuration(opts.waitTimeout))
	if opts.waitPattern != "" {
		fmt.Fprintf(&b, "Set WaitPattern /%s/\n", opts.waitPattern)
	}
	b.WriteRune('\n')

	for _, cmd := range commands {
		fmt.Fprintf(&b, "Type %s\n", quote(cmd))
		b.WriteString("Sleep 500ms\n")
		b.WriteString("Enter\n")
		b.WriteString("Wait\n")
		fmt.Fprintf(&b, "Sleep %s\n", tapeDuration(opts.pause))
		b.WriteRune('\n')
	}

	return b.String()
}
