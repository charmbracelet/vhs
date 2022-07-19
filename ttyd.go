package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func ttydCmd() *exec.Cmd {
	theme, _ := json.Marshal(DefaultTheme)

	return exec.Command(
		"ttyd",
		fmt.Sprintf("--port=%d", port),
		"-t", fmt.Sprintf("fontFamily='%s'", fontFamily),
		"-t", fmt.Sprintf("fontSize=%d", fontSize),
		"-t", fmt.Sprintf("lineHeight=%f", lineHeight),
		"-t", fmt.Sprintf("theme=%s", string(theme)),
		"-t", "customGlyphs=true",
		"zsh",
	)
}
