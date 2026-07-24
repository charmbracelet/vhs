package main

import (
	"strings"
	"sync"
	"testing"
)

// TestStart_NestedRecordingRejected ensures that starting a VHS recording
// while already inside a recording session (as signaled by recordingEnvVar,
// which VHS exports into the shell it drives) fails fast with a clear error
// instead of going on to spin up a second ttyd/browser pair, which is what
// used to cause a panic. See charmbracelet/vhs#761.
func TestStart_NestedRecordingRejected(t *testing.T) {
	t.Setenv(recordingEnvVar, "1")

	opts := DefaultVHSOptions()
	vhs := VHS{
		Options: &opts,
		mutex:   &sync.Mutex{},
	}

	err := vhs.Start()
	if err == nil {
		t.Fatal("expected an error when starting a nested recording, got nil")
	}
	if !strings.Contains(err.Error(), "already recording") {
		t.Fatalf("expected error to mention that vhs is already recording, got: %v", err)
	}
	if vhs.started {
		t.Fatal("vhs.started should remain false when a nested recording is rejected")
	}
}
