package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// sleepThreshold is the time at which if there has been no activity in the
// tape file we insert a Sleep command
const sleepThreshold = 500 * time.Millisecond

// EscapeSequences is a map of escape sequences to their VHS commands.
var EscapeSequences = map[string]string{
	"\x1b[A":  UP,
	"\x1b[B":  DOWN,
	"\x1b[C":  RIGHT,
	"\x1b[D":  LEFT,
	"\x1b[1~": HOME,
	"\x1b[2~": INSERT,
	"\x1b[3~": DELETE,
	"\x1b[4~": END,
	"\x1b[5~": PAGEUP,
	"\x1b[6~": PAGEDOWN,
	"\x01":    CTRL + "+A",
	"\x02":    CTRL + "+B",
	"\x03":    CTRL + "+C",
	"\x04":    CTRL + "+D",
	"\x05":    CTRL + "+E",
	"\x06":    CTRL + "+F",
	"\x07":    CTRL + "+G",
	"\x08":    BACKSPACE,
	"\x09":    TAB,
	"\x0b":    CTRL + "+K",
	"\x0c":    CTRL + "+L",
	"\x0d":    ENTER,
	"\x0e":    CTRL + "+N",
	"\x0f":    CTRL + "+O",
	"\x10":    CTRL + "+P",
	"\x11":    CTRL + "+Q",
	"\x12":    CTRL + "+R",
	"\x13":    CTRL + "+S",
	"\x14":    CTRL + "+T",
	"\x15":    CTRL + "+U",
	"\x16":    CTRL + "+V",
	"\x17":    CTRL + "+W",
	"\x18":    CTRL + "+X",
	"\x19":    CTRL + "+Y",
	"\x1a":    CTRL + "+Z",
	"\x1b":    ESCAPE,
	"\x7f":    BACKSPACE,
}

// Record is a command that starts a pseudo-terminal for the user to begin
// writing to, it records all the key presses on stdin and uses them to write
// Tape commands.
//
// vhs record > file.tape
func Record(_ *cobra.Command, _ []string) error {
	command := exec.Command(shell)

	terminal, err := pty.Start(command)
	if err != nil {
		return err
	}

	if err := pty.InheritSize(os.Stdin, terminal); err != nil {
		log.Printf("error resizing pty: %s", err)
	}

	prevState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	// We'll need to display the stdin on the screen but we'll also need a copy to
	// analyze later and create a tape file.
	tape := &bytes.Buffer{}
	in := io.MultiWriter(tape, terminal)

	go func() {
		var length int
		for {
			length = tape.Len()
			time.Sleep(sleepThreshold)
			if length == tape.Len() {
				// Tape has not changed in a while, write a Sleep command.
				tape.WriteString(fmt.Sprintf("\n%s\n", SLEEP))
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

	log.Println(inputToTape(tape.String()))
	return nil
}

var (
	cursorResponse = regexp.MustCompile(`\x1b\[\d+;\d+R`)
	oscResponse    = regexp.MustCompile(`\x1b\]\d+;rgb:....\/....\/....(\x07|\x1b\\)`)
)

// inputToTape takes input from a PTY stdin and converts it into a tape file.
func inputToTape(input string) string {
	// If the user exited the shell by typing exit don't record this in the
	// command.
	//
	// NOTE: this is not very robust as if someone types exii<BS>t it will not work
	// correctly and the exit will show up. In this case, the user should edit the
	// tape file.
	s := strings.TrimSuffix(strings.TrimSpace(input), "exit")

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
		}
		i += repeat - 1

		// We've encountered some non-command, assume that we need to type these
		// characters.
		if TokenType(lines[i]) == SLEEP {
			sleep := sleepThreshold * time.Duration(repeat)
			sanitized.WriteString(fmt.Sprintf("%s %s", TokenType(SLEEP), sleep))
		} else if strings.HasPrefix(lines[i], CTRL) {
			for j := 0; j < repeat; j++ {
				sanitized.WriteString("Ctrl" + strings.TrimPrefix(lines[i], CTRL) + "\n")
			}
			continue
		} else if strings.HasPrefix(lines[i], ALT) {
			for j := 0; j < repeat; j++ {
				sanitized.WriteString("Alt" + strings.TrimPrefix(lines[i], ALT) + "\n")
			}
			continue
		} else if IsCommand(TokenType(lines[i])) {
			sanitized.WriteString(fmt.Sprint(TokenType(lines[i])))
			if repeat > 1 {
				sanitized.WriteString(fmt.Sprint(" ", repeat))
			}
		} else {
			sanitized.WriteString(fmt.Sprintln(TokenType(TYPE), quote(lines[i])))
			continue
		}
		sanitized.WriteRune('\n')
	}

	return sanitized.String()
}

// quote wraps a string in (single or double) quotes
func quote(s string) string {
	if strings.ContainsRune(s, '"') && strings.ContainsRune(s, '\'') {
		return fmt.Sprintf("`%s`", s)
	}
	if strings.ContainsRune(s, '"') {
		return fmt.Sprintf(`'%s'`, s)
	}
	return fmt.Sprintf(`"%s"`, s)
}
