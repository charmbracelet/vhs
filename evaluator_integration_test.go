//go:build integration

package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Requires stock ttyd, Chrome or Chromium, and FFmpeg.
func TestEvaluateCreatesMedia(t *testing.T) {
	dir := t.TempDir()
	tape := fmt.Sprintf(`Output "%s/out.gif"
Set Shell bash
Set Width 400
Set Height 300
Set Framerate 10
Sleep 500ms
Screenshot "%s/out.png"
Sleep 500ms
`, dir, dir)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if errs := Evaluate(ctx, tape, io.Discard); len(errs) != 0 {
		t.Fatal(errs)
	}
	for _, name := range []string{"out.gif", "out.png"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		requireNoErr(t, err)
		_, _, err = image.Decode(bytes.NewReader(data))
		requireNoErr(t, err)
	}
}
