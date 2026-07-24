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

// recordingEnvVar is set in the environment of the shell that VHS drives so
// that any `vhs` invocation typed into that shell (e.g. a tape that itself
// runs `vhs`) can detect that it is nested inside an existing recording
// session. See VHS.Start for where this is checked.
const recordingEnvVar = "VHS_RECORDING"

// randomPort returns a random port number that is not in use.
func randomPort() int {
	addr, _ := net.Listen("tcp", ":0") //nolint:gosec,noctx
	_ = addr.Close()
	return addr.Addr().(*net.TCPAddr).Port
}

// buildTtyCmd builds the ttyd exec.Command on the given port.
func buildTtyCmd(port int, shell Shell) *exec.Cmd {
	args := []string{ //nolint:prealloc
		fmt.Sprintf("--port=%d", port),
		"--interface", "127.0.0.1",
		"-t", "rendererType=canvas",
		"-t", "disableResizeOverlay=true",
		"-t", "enableSixel=true",
		"-t", "customGlyphs=true",
		"--once", // will allow one connection and exit
		"--writable",
	}

	args = append(args, shell.Command...)

	cmd := exec.Command("ttyd", args...)
	env := append(shell.Env, os.Environ()...)
	cmd.Env = append(env, recordingEnvVar+"=1")
	return cmd
}
