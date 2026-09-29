// Package vhs tty.go spawns the ttyd process.
// It runs on the specified port and is generally meant to run in the background
// so that other processes (go-rod) can connect to the tty.
//
// xterm.js is used for rendering the terminal and can be adjusted using the Set command.
//
// Set FontFamily "DejaVu Sans Mono"
// Set FontSize 12
// Set Padding 50
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

// randomPort returns a random port number that is not in use.
func randomPort() int {
	addr, _ := net.Listen("tcp", ":0") //nolint:gosec,noctx
	_ = addr.Close()
	return addr.Addr().(*net.TCPAddr).Port
}

// terminalRenderer returns the xterm.js renderer type to request from ttyd,
// controlled by the VHS_RENDERER environment variable. It defaults to
// "canvas"; "dom" may be used on systems where canvas-based capture produces
// blank frames, such as headless or remote Windows environments.
func terminalRenderer() string {
	switch r := os.Getenv("VHS_RENDERER"); r {
	case "", "canvas":
		return "canvas"
	case "dom":
		return "dom"
	default:
		log.Printf("unknown VHS_RENDERER %q, using canvas renderer", r)
		return "canvas"
	}
}

// buildTtyCmd builds the ttyd exec.Command on the given port.
func buildTtyCmd(ctx context.Context, port int, shell Shell) *exec.Cmd {
	args := []string{ //nolint:prealloc
		fmt.Sprintf("--port=%d", port),
		"--interface", "127.0.0.1",
		"-t", "rendererType=" + terminalRenderer(),
		"-t", "disableResizeOverlay=true",
		"-t", "enableSixel=true",
		"-t", "customGlyphs=true",
		"--once", // will allow one connection and exit
		"--writable",
	}

	args = append(args, shell.Command...)

	cmd := exec.CommandContext(ctx, "ttyd", args...)
	if shell.Env != nil {
		cmd.Env = append(shell.Env, os.Environ()...)
	}
	return cmd
}
