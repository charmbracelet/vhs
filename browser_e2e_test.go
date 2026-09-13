package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

// TestStartBrowser launches a real browser through startBrowser, connects to
// it, opens a page, and shuts everything down again.
//
// It requires a browser; if none is installed on the machine one is
// downloaded, which is why it is opt-in: set VHS_TEST_BROWSER=1 to run it.
func TestStartBrowser(t *testing.T) {
	if os.Getenv("VHS_TEST_BROWSER") == "" {
		t.Skip("set VHS_TEST_BROWSER=1 to run this test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	browser, closer, err := startBrowser(ctx)
	if err != nil {
		t.Fatalf("startBrowser failed: %v", err)
	}

	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		t.Errorf("could not open page: %v", err)
	}
	if err := page.WaitLoad(); err != nil {
		t.Errorf("page did not load: %v", err)
	}

	res, err := page.Eval("() => 1 + 1")
	if err != nil {
		t.Errorf("could not evaluate: %v", err)
	} else if got := res.Value.Int(); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}

	if err := closer(); err != nil {
		t.Errorf("closer failed: %v", err)
	}
}

// TestRenderWritesOutput evaluates a minimal tape and checks that a file is
// actually written.
//
// Nothing else asserts on an output file, which is how a cancelled context
// reaching Render went unnoticed: frames rendered, ffmpeg was never run, and
// vhs exited 0.
//
// It requires a browser and ffmpeg; set VHS_TEST_BROWSER=1 to run it.
func TestRenderWritesOutput(t *testing.T) {
	if os.Getenv("VHS_TEST_BROWSER") == "" {
		t.Skip("set VHS_TEST_BROWSER=1 to run this test")
	}

	// The parser rejects an absolute path after Output, so the test runs from
	// the temporary directory and asks for a relative one.
	dir := t.TempDir()
	t.Chdir(dir)
	out := filepath.Join(dir, "out.gif")
	tape := "Output out.gif\nSet Width 400\nSet Height 200\nType \"a\"\nSleep 500ms\n"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if errs := Evaluate(ctx, tape, io.Discard); len(errs) > 0 {
		t.Fatalf("evaluate failed: %v", errs)
	}

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("no output written: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("output is empty")
	}
}
