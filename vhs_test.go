package main

import (
	"sync"
	"testing"
)

// TestTerminateIdempotent guards the fix for issue #738 — Evaluate()
// defers terminate() on every early return, and Record() also calls
// terminate() once context is done. Both paths run on the normal happy
// path, so terminate() must be safe to call twice.
func TestTerminateIdempotent(t *testing.T) {
	v := VHS{
		mutex:      &sync.Mutex{},
		started:    true,
		terminated: true, // simulate post-Record state: terminate() already ran
	}

	// Second call must be a no-op and must not deref nil browser/tty fields.
	if err := v.terminate(); err != nil {
		t.Errorf("second terminate() returned error: %v", err)
	}
}

// TestTerminateBeforeStartIsNoOp guards the fix for issue #738 — if
// Evaluate() reaches the deferred terminate() without Start() having
// succeeded (e.g. parser error path), terminate() must not panic
// dereferencing nil browser / tty fields.
func TestTerminateBeforeStartIsNoOp(t *testing.T) {
	v := VHS{
		mutex:   &sync.Mutex{},
		started: false,
	}
	if err := v.terminate(); err != nil {
		t.Errorf("terminate() before Start() returned error: %v", err)
	}
}
