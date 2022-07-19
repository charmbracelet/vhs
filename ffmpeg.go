package main

import (
	"fmt"
	"os/exec"
)

const maxColors = 256
const scale = "scale=1200:600:-1"
const flags = "lanczos"

func ffmpegCmd() *exec.Cmd {
	return exec.Command(
		"ffmpeg",
		"-i", capturesPath,
		"-vf", fmt.Sprintf("scale=%s:flags=%s,split[s0][s1];[s0]palettegen=max_colors=%d[p];[s1][p]paletteuse", scale, flags, maxColors),
		gifPath,
	)
}
