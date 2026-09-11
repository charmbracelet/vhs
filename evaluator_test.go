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

func TestEvaluateGIF(t *testing.T) {
	if os.Getenv("VHS_TEST_BROWSER") == "" {
		t.Skip("set VHS_TEST_BROWSER=1 to run this test")
	}
	for _, name := range []string{"ttyd", "ffmpeg"} {
		if _, err := exec.LookPath(name); err != nil {
			t.Fatalf("%s is required: %v", name, err)
		}
	}

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
