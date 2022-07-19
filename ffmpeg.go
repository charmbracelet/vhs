package main

import (
	"fmt"
	"os/exec"
)

func ffmpegCmd() *exec.Cmd {
	return exec.Command(
		"ffmpeg",
		"-y",
		"-i", capturesPath,
		"-framerate", fmt.Sprint(framerate),
		"-vf", fmt.Sprintf("fps=%d,scale=%d:-1:flags=%s,split[s0][s1];[s0]palettegen=max_colors=%d[p];[s1][p]paletteuse", framerate, width, ffmpegFlags, maxColors),
		gifPath,
	)
}
