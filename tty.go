package dolly

import (
	"fmt"
	"os/exec"
)

// StartTTY starts the ttyd process on the given port and options.
func StartTTY(port int) *exec.Cmd {
	cmd := exec.Command(
		"ttyd", fmt.Sprintf("--port=%d", port),
		"-t", "customGlyphs=true",
		"zsh",
	)
	return cmd
}
