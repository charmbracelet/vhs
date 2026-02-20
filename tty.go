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
	"fmt"
	"net"
	"os"
	"os/exec"
)

// randomPort returns a random port number that is not in use.
func randomPort() int {
	addr, _ := net.Listen("tcp", ":0") //nolint:gosec
	_ = addr.Close()
	return addr.Addr().(*net.TCPAddr).Port
}

// buildTtyCmd builds the ttyd exec.Command on the given port.
func buildTtyCmd(port int, shell Shell, promptColor, prompt string) *exec.Cmd {
	args := []string{
		fmt.Sprintf("--port=%d", port),
		"--interface", "127.0.0.1",
		"-t", "rendererType=canvas",
		"-t", "disableResizeOverlay=true",
		"-t", "enableSixel=true",
		"-t", "customGlyphs=true",
		"--once", // will allow one connection and exit
		"--writable",
	}

	env, command := ShellConfig(shell.Name, promptColor, prompt)
	args = append(args, command...)

	cmd := exec.Command("ttyd", args...)
	if env != nil {
		cmd.Env = append(env, os.Environ()...)
	}
	return cmd
}
