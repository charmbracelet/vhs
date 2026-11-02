package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
	"github.com/charmbracelet/x/xpty"
)

// Terminal is a shell running on a pseudo-terminal whose output is fed into a
// virtual terminal emulator.
type Terminal struct {
	pty xpty.Pty
	emu *vt.SafeEmulator
	cmd *exec.Cmd

	mu sync.Mutex

	closeOnce sync.Once
	closeErr  error
	// scrollTop is the scrollback line at the top of the viewport while the
	// user has scrolled up, so new output doesn't move what is being viewed.
	scrollTop int
	scrolled  bool
	// lastInput is when input was last sent, used to restart cursor blink.
	lastInput time.Time
}

// StartTerminal starts the given shell on a new pseudo-terminal of the given
// size.
func StartTerminal(ctx context.Context, shell Shell, cols, rows int) (*Terminal, error) {
	if len(shell.Command) == 0 {
		return nil, errors.New("no shell command")
	}

	pty, err := xpty.NewPty(cols, rows)
	if err != nil {
		return nil, fmt.Errorf("could not create pty: %w", err)
	}

	cmd := exec.CommandContext(ctx, shell.Command[0], shell.Command[1:]...) //nolint:gosec
	cmd.Env = append(os.Environ(), "TERM=xterm-256color", "COLORTERM=truecolor")
	cmd.Env = append(cmd.Env, shell.Env...)
	if err := pty.Start(cmd); err != nil {
		_ = pty.Close()
		return nil, fmt.Errorf("could not start shell %s: %w", shell.Command[0], err)
	}

	emu := vt.NewSafeEmulator(cols, rows)
	t := &Terminal{pty: pty, emu: emu, cmd: cmd, lastInput: time.Now()}

	// Shell output to the emulator, and emulator replies and input to the
	// shell.
	go io.Copy(emu, pty) //nolint:errcheck
	go io.Copy(pty, emu) //nolint:errcheck

	return t, nil
}

// SetTheme sets the default and ANSI colors of the terminal.
func (t *Terminal) SetTheme(theme Theme) {
	fg, bg, cursor, _ := theme.colors()
	t.emu.SetDefaultForegroundColor(fg)
	t.emu.SetDefaultBackgroundColor(bg)
	t.emu.SetDefaultCursorColor(cursor)
	for i, c := range theme.palette() {
		t.emu.SetIndexedColor(i, c)
	}
}

// Resize resizes the terminal and its pty.
func (t *Terminal) Resize(cols, rows int) error {
	t.emu.Resize(cols, rows)
	if err := t.pty.Resize(cols, rows); err != nil {
		return fmt.Errorf("could not resize pty: %w", err)
	}
	return nil
}

// Size returns the number of columns and rows of the terminal.
func (t *Terminal) Size() (cols, rows int) {
	return t.emu.Width(), t.emu.Height()
}

// userInput scrolls back to the bottom and restarts the cursor blink, like
// xterm.js does when the user types.
func (t *Terminal) userInput() {
	t.mu.Lock()
	t.scrolled = false
	t.lastInput = time.Now()
	t.mu.Unlock()
}

// SendKey sends a key press to the shell, encoded according to the current
// terminal modes.
func (t *Terminal) SendKey(k uv.KeyPressEvent) {
	t.userInput()
	t.emu.SendKey(k)
}

// SendText sends text to the shell as if typed.
func (t *Terminal) SendText(s string) {
	t.userInput()
	t.emu.SendText(s)
}

// LastInput returns when input was last sent.
func (t *Terminal) LastInput() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastInput
}

// Scroll moves the viewport by delta lines: negative scrolls up into the
// scrollback, positive scrolls back down.
func (t *Terminal) Scroll(delta int) {
	if t.emu.IsAltScreen() {
		return
	}
	sb := t.emu.ScrollbackLen()

	t.mu.Lock()
	defer t.mu.Unlock()
	offset := min(max(t.viewOffset(sb)-delta, 0), sb)
	t.scrolled = offset > 0
	t.scrollTop = sb - offset
}

// ViewOffset returns how many lines the viewport is scrolled up from the
// bottom.
func (t *Terminal) ViewOffset() int {
	if t.emu.IsAltScreen() {
		return 0
	}
	sb := t.emu.ScrollbackLen()
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.viewOffset(sb)
}

func (t *Terminal) viewOffset(sb int) int {
	if !t.scrolled {
		return 0
	}
	return min(max(sb-t.scrollTop, 0), sb)
}

// Snapshot returns a copy of the visible viewport, including scrollback lines
// when scrolled up.
func (t *Terminal) Snapshot() uv.ScreenBuffer {
	cols, rows := t.Size()
	buf := uv.NewScreenBuffer(cols, rows)
	t.emu.Draw(buf, buf.Bounds())

	offset := t.ViewOffset()
	if offset == 0 {
		return buf
	}

	// Rows are filled bottom-up: screen rows move down by offset, so each
	// source row is read before it is overwritten.
	sb := t.emu.ScrollbackLen()
	for y := rows - 1; y >= 0; y-- {
		line := sb - offset + y
		for x := 0; x < cols; {
			var c *uv.Cell
			if line >= sb {
				c = buf.CellAt(x, line-sb)
			} else {
				c = t.emu.ScrollbackCellAt(x, line)
			}
			if c == nil || c.IsZero() {
				c = &uv.EmptyCell
			}
			buf.SetCell(x, y, c)
			x += max(c.Width, 1)
		}
	}
	return buf
}

// Lines returns the text of the visible viewport, one string per row, with
// trailing spaces removed.
func (t *Terminal) Lines() []string {
	buf := t.Snapshot()
	lines := make([]string, buf.Height())
	for y := range lines {
		lines[y] = strings.TrimRight(buf.Line(y).String(), " ")
	}
	return lines
}

// CurrentLine returns the text of the line the cursor is on.
func (t *Terminal) CurrentLine() string {
	cols, rows := t.Size()
	buf := uv.NewScreenBuffer(cols, rows)
	t.emu.Draw(buf, buf.Bounds())
	y := t.emu.CursorPosition().Y
	if y < 0 || y >= rows {
		return ""
	}
	return strings.TrimRight(buf.Line(y).String(), " ")
}

// Colors returns the current default foreground and background colors,
// which programs can change with OSC 10 and 11.
func (t *Terminal) Colors() (fg, bg color.Color) {
	return t.emu.ForegroundColor(), t.emu.BackgroundColor()
}

// Close stops the shell and releases the pty. It is safe to call more than
// once.
func (t *Terminal) Close() error {
	t.closeOnce.Do(func() {
		_ = t.emu.Close()
		if t.cmd.Process != nil {
			_ = t.cmd.Process.Kill()
			_ = xpty.WaitProcess(context.Background(), t.cmd)
		}
		if err := t.pty.Close(); err != nil {
			t.closeErr = fmt.Errorf("could not close pty: %w", err)
		}
	})
	return t.closeErr
}
