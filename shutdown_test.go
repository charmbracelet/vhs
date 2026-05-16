package main

import (
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestShutdownEmptyVHS(t *testing.T) {
	// shutdown should be a no-op (and never panic) on a fresh VHS that
	// hasn't been Start()ed: both tty and browser are nil.
	var vhs VHS
	if err := vhs.shutdown(); err != nil {
		t.Fatalf("shutdown on empty VHS: %v", err)
	}
}

func TestShutdownKillsTTY(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sleep binary not available on Windows runners")
	}

	// Use a long-running sleep as a stand-in for ttyd so we can check that
	// shutdown actually killed it. Regression for the ttyd leak when
	// Evaluate returns before Record begins.
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start sleep: %v", err)
	}

	vhs := VHS{tty: cmd}
	if err := vhs.shutdown(); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
		// The Wait return value is intentionally ignored: the process was
		// killed, so an error is expected.
	case <-time.After(2 * time.Second):
		t.Fatal("ttyd stand-in still alive 2s after shutdown")
	}
}

func TestShutdownIsIdempotent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sleep binary not available on Windows runners")
	}

	// terminate() already runs at the end of the happy-path recording, then
	// the deferred close() runs from Evaluate. Calling shutdown twice in a
	// row should not surface an error from the second call.
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start sleep: %v", err)
	}
	vhs := VHS{tty: cmd}
	if err := vhs.shutdown(); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := vhs.shutdown(); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}
}
