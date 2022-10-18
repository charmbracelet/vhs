// tty is what spawns the ttyd process.
// It runs on the specified port and is generally meant to run in the background
// so that other processes (go-rod) can connect to the tty.
//
// xterm.js is used for rendering the terminal and can be adjusted using the Set command.
//
// Set FontFamily SF Mono
// Set FontSize 12
// Set Padding 1em
package vhs

import (
	"fmt"
	"net"
	"os/exec"
)

// randomPort returns a random port number that is not in use.
func randomPort() int {
	addr, _ := net.Listen("tcp", ":0")
	addr.Close()
	return addr.Addr().(*net.TCPAddr).Port
}

// StartTTY starts the ttyd process on the given port.
func StartTTY(port int) *exec.Cmd {
	cmd := exec.Command(
		"ttyd", fmt.Sprintf("--port=%d", port),
		"-t", "rendererType=dom",
		"-t", "disableResizeOverlay=true",
		"-t", "cursorBlink=true",
		"-t", "customGlyphs=true",
		"bash", "--login",
	)
	return cmd
}
