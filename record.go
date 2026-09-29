package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/vhs/token"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// sleepThreshold is the time at which if there has been no activity in the
// tape file we insert a Sleep command.
const sleepThreshold = 500 * time.Millisecond

// EscapeSequences is a map of escape sequences to their VHS commands.
var EscapeSequences = map[string]string{
	"\x1b[A":  token.UP,
	"\x1b[B":  token.DOWN,
	"\x1b[C":  token.RIGHT,
	"\x1b[D":  token.LEFT,
	"\x1b[1~": token.HOME,
	"\x1b[2~": token.INSERT,
	"\x1b[3~": token.DELETE,
	"\x1b[4~": token.END,
	"\x1b[5~": token.PAGE_UP,
	"\x1b[6~": token.PAGE_DOWN,
	"\x01":    token.CTRL + "+A",
	"\x02":    token.CTRL + "+B",
	"\x03":    token.CTRL + "+C",
	"\x04":    token.CTRL + "+D",
	"\x05":    token.CTRL + "+E",
	"\x06":    token.CTRL + "+F",
	"\x07":    token.CTRL + "+G",
	"\x08":    token.BACKSPACE,
	"\x09":    token.TAB,
	"\x0b":    token.CTRL + "+K",
	"\x0c":    token.CTRL + "+L",
	"\x0d":    token.ENTER,
	"\x0e":    token.CTRL + "+N",
	"\x0f":    token.CTRL + "+O",
	"\x10":    token.CTRL + "+P",
	"\x11":    token.CTRL + "+Q",
	"\x12":    token.CTRL + "+R",
	"\x13":    token.CTRL + "+S",
	"\x14":    token.CTRL + "+T",
	"\x15":    token.CTRL + "+U",
	"\x16":    token.CTRL + "+V",
	"\x17":    token.CTRL + "+W",
	"\x18":    token.CTRL + "+X",
	"\x19":    token.CTRL + "+Y",
	"\x1a":    token.CTRL + "+Z",
	"\x1b":    token.ESCAPE,
	"\x7f":    token.BACKSPACE,
}

// Record is a command that starts a pseudo-terminal for the user to begin
// writing to, it records all the key presses on stdin and uses them to write
// Tape commands.
//
//	vhs record > file.tape
//
//nolint:wrapcheck
func Record(cmd *cobra.Command, _ []string) error {
	terminal, err := startRecordingTerminal(cmd.Context(), shell)
	if err != nil {
		return err
	}

	prevState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	// We'll need to display the stdin on the screen but we'll also need a copy to
	// analyze later and create a tape file.
	tape := &bytes.Buffer{}
	in := io.MultiWriter(tape, terminal)

	if shell != defaultShell {
		_, _ = fmt.Fprintf(tape, "%s Shell %s\n", token.SET, shell)
	}

	go func() {
		var length int
		for {
			length = tape.Len()
			time.Sleep(sleepThreshold)
			if length == tape.Len() {
				// Tape has not changed in a while, write a Sleep command.
				_, _ = fmt.Fprintf(tape, "\n%s\n", token.SLEEP)
			}
		}
	}()

	// Write to the buffer and PTY's stdin and stderr so that stdout is reserved
	// for the output tape file.
	go func() { _, _ = io.Copy(in, os.Stdin) }()
	_, _ = io.Copy(os.Stderr, terminal)

	// PTY cleanup and restore terminal
	_ = terminal.Close()
	_ = term.Restore(int(os.Stdin.Fd()), prevState)

	fmt.Println(inputToTape(tape.String()))
	return nil
}

var (
	cursorResponse = regexp.MustCompile(`\x1b\[\d+;\d+R`)
	oscResponse    = regexp.MustCompile(`\x1b\]\d+;rgb:....\/....\/....(\x07|\x1b\\)`)
	// win32InputEvent matches a Windows console input-mode key event
	// (Vk;Sc;Uc;Kd;Cs;Rc) as emitted by terminals in win32-input-mode.
	win32InputEvent = regexp.MustCompile(`\x1b\[(\d+);(\d+);(\d+);(\d+);(\d+);(\d+)_`)
	// focusEvent matches the focus in/out reports terminals emit once
	// focus reporting is enabled.
	focusEvent = regexp.MustCompile(`\x1b\[[IO]`)
)

// win32VirtualKeys maps Windows virtual-key codes to the byte sequences a
// Unix pseudo-terminal emits for the same keys.
var win32VirtualKeys = map[int]string{
	0x08: "\x08",    // backspace
	0x09: "\x09",    // tab
	0x0d: "\x0d",    // enter
	0x1b: "\x1b",    // escape
	0x21: "\x1b[5~", // page up
	0x22: "\x1b[6~", // page down
	0x23: "\x1b[4~", // end
	0x24: "\x1b[1~", // home
	0x25: "\x1b[D",  // left
	0x26: "\x1b[A",  // up
	0x27: "\x1b[C",  // right
	0x28: "\x1b[B",  // down
	0x2d: "\x1b[2~", // insert
	0x2e: "\x1b[3~", // delete
}

// translateConsoleInput rewrites console input-mode event sequences into the
// byte stream a Unix pseudo-terminal would produce for the same keystrokes:
// win32-input-mode key events become their key equivalents and focus reports
// are dropped. Key releases, modifier presses, and events that produce no
// text have no Unix counterpart and are dropped as well.
func translateConsoleInput(s string) string {
	s = focusEvent.ReplaceAllString(s, "")
	return win32InputEvent.ReplaceAllStringFunc(s, func(seq string) string {
		m := win32InputEvent.FindStringSubmatch(seq)
		vk, _ := strconv.Atoi(m[1])
		uc, _ := strconv.Atoi(m[3])
		kd, _ := strconv.Atoi(m[4])
		cs, _ := strconv.Atoi(m[5])
		if kd == 0 {
			return "" // key release
		}
		if b, ok := win32VirtualKeys[vk]; ok {
			return b
		}
		switch vk {
		case 0x10, 0x11, 0x12, 0x5b, 0x5c:
			return "" // shift, ctrl, alt, and win keys produce no input
		}
		if uc == 0 {
			return "" // event produced no text
		}
		if cs&0x0c != 0 && uc < 0x80 {
			return string(byte(uc & 0x1f)) // ctrl+key -> control byte
		}
		return string(rune(uc))
	})
}

// inputToTape takes input from a PTY stdin and converts it into a tape file.
func inputToTape(input string) string {
	// If the user exited the shell by typing exit don't record this in the
	// command.
	//
	// NOTE: this is not very robust as if someone types exii<BS>t it will not work
	// correctly and the exit will show up. In this case, the user should edit the
	// tape file.
	s := strings.TrimSuffix(strings.TrimSpace(translateConsoleInput(input)), "exit")

	// Remove cursor / osc responses
	s = cursorResponse.ReplaceAllString(s, "")
	s = oscResponse.ReplaceAllString(s, "")

	// Substitute escape sequences for commands
	for sequence, command := range EscapeSequences {
		s = strings.ReplaceAll(s, sequence, "\n"+command+"\n")
	}

	s = strings.ReplaceAll(s, "\n\n", "\n")

	var sanitized strings.Builder
	lines := strings.Split(s, "\n")

	for i := 0; i < len(lines)-1; i++ {
		// Group repeated commands to compress file and make it more readable.
		repeat := 1
		for lines[i] == lines[i+repeat] {
			repeat++
			if i+repeat == len(lines) {
				break
			}
		}
		i += repeat - 1

		// We've encountered some non-command, assume that we need to type these
		// characters.
		if token.Type(lines[i]) == token.SLEEP { //nolint:nestif
			sleep := sleepThreshold * time.Duration(repeat)
			if sleep >= time.Minute {
				_, _ = fmt.Fprintf(&sanitized, "%s %gs", token.Type(token.SLEEP), sleep.Seconds())
			} else {
				_, _ = fmt.Fprintf(&sanitized, "%s %s", token.Type(token.SLEEP), sleep)
			}
		} else if strings.HasPrefix(lines[i], token.CTRL) {
			for j := 0; j < repeat; j++ {
				sanitized.WriteString("Ctrl" + strings.TrimPrefix(lines[i], token.CTRL) + "\n")
			}
			continue
		} else if strings.HasPrefix(lines[i], token.ALT) {
			for j := 0; j < repeat; j++ {
				sanitized.WriteString("Alt" + strings.TrimPrefix(lines[i], token.ALT) + "\n")
			}
			continue
		} else if strings.HasPrefix(lines[i], token.SET) {
			sanitized.WriteString("Set" + strings.TrimPrefix(lines[i], token.SET))
		} else if token.IsCommand(token.Type(lines[i])) {
			_, _ = fmt.Fprint(&sanitized, token.Type(lines[i]))
			if repeat > 1 {
				_, _ = fmt.Fprint(&sanitized, " ", repeat)
			}
		} else {
			if lines[i] != "" {
				_, _ = fmt.Fprintln(&sanitized, token.Type(token.TYPE), quote(lines[i]))
			}
			continue
		}
		sanitized.WriteRune('\n')
	}

	return sanitized.String()
}

// quote wraps a string in (single or double) quotes.
func quote(s string) string {
	if strings.ContainsRune(s, '"') && strings.ContainsRune(s, '\'') {
		return fmt.Sprintf("`%s`", s)
	}
	if strings.ContainsRune(s, '"') {
		return fmt.Sprintf(`'%s'`, s)
	}
	return fmt.Sprintf(`"%s"`, s)
}
