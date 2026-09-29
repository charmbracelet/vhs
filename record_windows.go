//go:build windows

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/x/conpty"
	"golang.org/x/sys/windows"
)

// conptyTerminal adapts a Windows pseudo-console to the io.ReadWriteCloser
// interface used by the Unix PTY implementation.
type conptyTerminal struct {
	cpty   *conpty.ConPty
	handle windows.Handle
	exited chan struct{}
	once   sync.Once
}

// startRecordingTerminal starts the user's shell on a fresh Windows
// pseudo-console (ConPTY) and returns it as a stream: reads yield the shell's
// output, writes feed its input.
func startRecordingTerminal(ctx context.Context, shell string) (*conptyTerminal, error) {
	command := exec.CommandContext(ctx, shell)
	command.Env = append(os.Environ(), "VHS_RECORD=true")

	w, h := conpty.DefaultWidth, conpty.DefaultHeight
	if cols, rows, err := consoleSize(); err == nil {
		w, h = cols, rows
	}

	cpty, err := conpty.New(w, h, 0)
	if err != nil {
		return nil, fmt.Errorf("create pseudo-console: %w", err)
	}

	_, handle, err := cpty.Spawn(command.Args[0], command.Args, &syscall.ProcAttr{
		Dir: command.Dir,
		Env: command.Env,
	})
	if err != nil {
		_ = cpty.Close()
		return nil, fmt.Errorf("spawn %q in pseudo-console: %w", command.Args[0], err)
	}

	t := &conptyTerminal{cpty: cpty, handle: windows.Handle(handle), exited: make(chan struct{})}

	// A console client exiting does not always close the output pipe
	// promptly; watch the process handle and tear the pseudo-console down so
	// blocked readers finish. The context canceling (e.g. vhs shutting down)
	// kills the shell.
	go func() {
		_, _ = windows.WaitForSingleObject(t.handle, windows.INFINITE)
		close(t.exited)
		t.once.Do(func() { _ = t.cpty.Close() })
	}()
	go func() {
		select {
		case <-ctx.Done():
			_ = windows.TerminateProcess(t.handle, 1)
		case <-t.exited:
		}
	}()

	return t, nil
}

// consoleSize reports the dimensions of the attached console's screen buffer.
func consoleSize() (cols, rows int, err error) {
	var info windows.ConsoleScreenBufferInfo
	for _, f := range []*os.File{os.Stdout, os.Stderr, os.Stdin} {
		if err = windows.GetConsoleScreenBufferInfo(windows.Handle(f.Fd()), &info); err == nil {
			return int(info.Window.Right - info.Window.Left + 1),
				int(info.Window.Bottom - info.Window.Top + 1), nil
		}
	}
	return 0, 0, fmt.Errorf("no console screen buffer: %w", err)
}

func (t *conptyTerminal) Read(p []byte) (int, error)  { return t.cpty.Read(p) }
func (t *conptyTerminal) Write(p []byte) (int, error) { return t.cpty.Write(p) }

// Close waits for the shell to exit, terminating it after a grace period if
// it is still running, then releases the pseudo-console and process handle.
func (t *conptyTerminal) Close() error {
	const grace = 5 * time.Second
	select {
	case <-t.exited:
	case <-time.After(grace):
		_ = windows.TerminateProcess(t.handle, 1)
		<-t.exited
	}
	t.once.Do(func() { _ = t.cpty.Close() })
	return windows.CloseHandle(t.handle)
}
