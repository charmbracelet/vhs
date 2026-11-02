//go:build !windows

package main

import (
	"context"
	"testing"
	"time"

	uv "github.com/charmbracelet/ultraviolet"
)

func startTestShell(t *testing.T, cols, rows int) *Terminal {
	t.Helper()
	term, err := StartTerminal(context.Background(), Shell{Command: []string{"sh"}, Env: []string{"PS1=> "}}, cols, rows)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = term.Close() })
	return term
}

func TestTerminalShell(t *testing.T) {
	term := startTestShell(t, 40, 10)
	term.SendText("echo hello\r")
	waitForLine(t, term, "hello")
	waitForLine(t, term, ">")
	if got := term.CurrentLine(); got != ">" {
		t.Fatalf("expected the prompt on the current line, got %q", got)
	}
}

func TestTerminalResize(t *testing.T) {
	term := startTestShell(t, 40, 10)
	if err := term.Resize(50, 12); err != nil {
		t.Fatal(err)
	}
	term.SendText("stty size\r")
	waitForLine(t, term, "12 50")
	if cols, rows := term.Size(); cols != 50 || rows != 12 {
		t.Fatalf("expected 50x12, got %dx%d", cols, rows)
	}
}

func TestTerminalCtrlC(t *testing.T) {
	// Ctrl+C only interrupts the foreground program when the shell runs with
	// the pty as its controlling terminal.
	term := startTestShell(t, 40, 10)
	term.SendText("sleep 10\r")
	time.Sleep(300 * time.Millisecond)
	k, _ := ctrlKey("c")
	term.SendKey(k)
	term.SendText("echo done\r")
	// Times out before sleep would end on its own.
	waitForLine(t, term, "done")
}

func TestTerminalKeysFollowModes(t *testing.T) {
	// cat -v shows the bytes the keys send.
	term := startTestShell(t, 40, 10)
	term.SendText("cat -v\r")
	time.Sleep(200 * time.Millisecond)
	term.SendKey(uv.KeyPressEvent{Code: keyLeft.Code, Mod: uv.ModCtrl})
	term.SendKey(keyEnter)
	waitForLine(t, term, "^[[1;5D")
}

func TestTerminalCloseTwice(t *testing.T) {
	term := startTestShell(t, 40, 10)
	if err := term.Close(); err != nil {
		t.Fatal(err)
	}
	if err := term.Close(); err != nil {
		t.Fatal(err)
	}
}
