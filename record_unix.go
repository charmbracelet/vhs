//go:build !windows
// +build !windows

package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// startRecordingTerminal starts the user's shell on a fresh pseudo-terminal and
// returns it as a stream: reads yield the shell's output, writes feed its
// input. The terminal tracks the size of os.Stdin for the life of the process.
func startRecordingTerminal(ctx context.Context, shell string) (io.ReadWriteCloser, error) {
	command := exec.CommandContext(ctx, shell)

	command.Env = append(os.Environ(), "VHS_RECORD=true")

	terminal, err := pty.Start(command)
	if err != nil {
		return nil, err
	}

	if err := pty.InheritSize(os.Stdin, terminal); err != nil {
		log.Printf("error resizing pty: %s", err)
	}

	return terminal, nil
}
