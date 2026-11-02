package main

import (
	"context"
	"errors"
	"fmt"
	imagegif "image/gif"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRenderMissingFFmpeg(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	v := New()
	v.totalFrames = 1
	t.Cleanup(func() { _ = v.Cleanup() })
	v.Options.Video.Output.GIF = filepath.Join(t.TempDir(), "output.gif")

	if err := v.Render(context.Background()); !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("expected missing ffmpeg error, got %v", err)
	}
}

// requireE2E skips end-to-end tests on machines without ffmpeg or fonts.
func requireE2E(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg is required to run this test")
	}
	if _, err := newFontSet(defaultFontFamily, defaultFontSize); err != nil {
		t.Skipf("fonts are required to run this test: %v", err)
	}
}

func TestEvaluateRowsColumns(t *testing.T) {
	requireE2E(t)
	dir := t.TempDir()
	gifPath := filepath.Join(dir, "output.gif")
	txtPath := filepath.Join(dir, "output.txt")
	tape := fmt.Sprintf(
		"Output %q\nOutput %q\nSet Shell bash\nSet FontSize 20\nSet Padding 10\nSet Rows 8\nSet Columns 30\nType \"stty size\"\nEnter\nSleep 500ms\n",
		filepath.ToSlash(gifPath), filepath.ToSlash(txtPath),
	)
	if errs := Evaluate(context.Background(), tape, io.Discard); len(errs) != 0 {
		t.Fatalf("Evaluate failed: %v", errs)
	}

	// The terminal has exactly the rows and columns asked for.
	out, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "\n8 30\n") {
		t.Fatalf("expected stty to report 8 rows and 30 columns, got:\n%s", out)
	}

	// The video fits that grid plus padding.
	opts := DefaultVHSOptions()
	opts.FontSize = 20
	cellW, cellH := newTestRenderer(t, opts).CellSize()
	file, err := os.Open(gifPath)
	if err != nil {
		t.Fatalf("output was not created: %v", err)
	}
	defer file.Close()
	cfg, err := imagegif.DecodeConfig(file)
	if err != nil {
		t.Fatalf("output is not a valid GIF: %v", err)
	}
	wantW, wantH := roundUpToEven(30*cellW+20), roundUpToEven(8*cellH+20)
	if cfg.Width != wantW || cfg.Height != wantH {
		t.Fatalf("expected a %dx%d GIF, got %dx%d", wantW, wantH, cfg.Width, cfg.Height)
	}
}

func TestEvaluateGIF(t *testing.T) {
	requireE2E(t)

	for _, cancelRender := range []bool{false, true} {
		t.Run(fmt.Sprintf("cancel_render=%t", cancelRender), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			output := filepath.Join(t.TempDir(), "output.gif")
			tape := fmt.Sprintf("Output %q\nSet Shell bash\nSet Width 320\nSet Height 200\nType hello\nSleep 500ms\n", filepath.ToSlash(output))
			errs := Evaluate(ctx, tape, io.Discard, func(_ *VHS) {
				if cancelRender {
					cancel()
				}
			})
			if cancelRender {
				if len(errs) != 1 || !errors.Is(errs[0], context.Canceled) {
					t.Fatalf("expected cancellation error, got %v", errs)
				}
				return
			}
			if len(errs) != 0 {
				t.Fatalf("Evaluate failed: %v", errs)
			}
			file, err := os.Open(output)
			if err != nil {
				t.Fatalf("output was not created: %v", err)
			}
			defer file.Close()
			if _, err := imagegif.DecodeAll(file); err != nil {
				t.Fatalf("output is not a valid GIF: %v", err)
			}
		})
	}
}
