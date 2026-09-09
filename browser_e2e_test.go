package main

import (
	"context"
	"os"
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
