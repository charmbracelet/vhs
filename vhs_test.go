package main

import "testing"

func TestTtyHostMatchesBrowserURL(t *testing.T) {
	const port = 19760

	cmd := buildTtyCmd(port, Shell{Command: []string{"sh"}})
	iface := ""
	for i, arg := range cmd.Args {
		if arg == "--interface" && i+1 < len(cmd.Args) {
			iface = cmd.Args[i+1]
			break
		}
	}

	if iface != ttyHost {
		t.Fatalf("expected ttyd to bind to %q, got %q", ttyHost, iface)
	}

	wantURL := "http://" + iface + ":19760"
	if gotURL := ttyURL(port); gotURL != wantURL {
		t.Fatalf("expected browser URL %q, got %q", wantURL, gotURL)
	}
}
