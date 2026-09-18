package main

import (
	"context"
	"path/filepath"
	"testing"
)

// TestRenderReportsEncodeFailure locks in two halves of the same regression.
//
// Render builds its encoders with exec.CommandContext, so a command started on
// a context that is already done is killed at Start and writes nothing. Before
// this was fixed, Evaluate handed Render the recording context, which teardown
// had just cancelled, and Render logged the empty output and returned nil. VHS
// printed "Creating out.gif...", exited 0, and produced no file
// (https://github.com/charmbracelet/vhs/issues/787).
//
// A cancelled context stands in here for any failing encode: whatever the
// cause, Render must report it rather than return nil.
func TestRenderReportsEncodeFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	v := New()
	v.Options.Video.Output.GIF = filepath.Join(t.TempDir(), "out.gif")
	v.totalFrames = 1

	if err := v.Render(ctx); err == nil {
		t.Fatal("Render returned nil for an encode that never ran; encode failures must reach the caller")
	}
}
