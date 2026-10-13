package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

func TestRenderGIFMemory(t *testing.T) {
	if os.Getenv("VHS_TEST_FFMPEG_MEMORY") == "" {
		t.Skip("set VHS_TEST_FFMPEG_MEMORY=1 to run this test")
	}
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("ffmpeg is required")
	}
	ffmpeg, err = filepath.Abs(ffmpeg)
	requireNoErr(t, err)
	testBinary, err := os.Executable()
	requireNoErr(t, err)
	t.Setenv("VHS_TEST_FFMPEG_BINARY", ffmpeg)
	t.Setenv("VHS_TEST_BINARY", testBinary)
	bin := t.TempDir()
	script := "#!/bin/sh\nexec \"$VHS_TEST_BINARY\" -test.run='^TestGIFEncoderHelper$' -- \"$@\"\n"
	requireNoErr(t, os.WriteFile(filepath.Join(bin, "ffmpeg"), []byte(script), 0o700))
	t.Setenv("PATH", bin)

	testRenderGIF(t, &StyleOptions{Width: 800, Height: 600, BackgroundColor: "#000000"}, 400, 399, 2)
}

func TestGIFEncoderHelper(t *testing.T) {
	ffmpeg := os.Getenv("VHS_TEST_FFMPEG_BINARY")
	if ffmpeg == "" {
		return
	}
	args := []string{ffmpeg, "-filter_complex_threads", "1", "-filter_threads", "1"}
	for i, arg := range os.Args {
		if arg != "--" {
			continue
		}
		for j, arg := range os.Args[i+1:] {
			if arg == "-i" || j == len(os.Args[i+1:])-1 {
				args = append(args, "-threads", "1")
			}
			args = append(args, arg)
		}
		break
	}
	requireNoErr(t, os.Setenv("MALLOC_ARENA_MAX", "2"))
	stack := &syscall.Rlimit{Cur: 1024 * 1024, Max: 1024 * 1024}
	requireNoErr(t, syscall.Setrlimit(syscall.RLIMIT_STACK, stack))
	limit := &syscall.Rlimit{Cur: 512 * 1024 * 1024, Max: 512 * 1024 * 1024}
	requireNoErr(t, syscall.Setrlimit(syscall.RLIMIT_AS, limit))
	requireNoErr(t, syscall.Exec(ffmpeg, args, os.Environ()))
}
