package main

import (
	"context"
	"testing"
)

func argsContain(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}

func TestBuildTtyCmdRenderer(t *testing.T) {
	shell := Shell{Command: []string{"bash"}}

	cmd := buildTtyCmd(context.Background(), 1976, shell)
	if !argsContain(cmd.Args, "rendererType=canvas") {
		t.Fatalf("expected canvas renderer by default, got %v", cmd.Args)
	}

	t.Setenv("VHS_RENDERER", "dom")
	cmd = buildTtyCmd(context.Background(), 1976, shell)
	if !argsContain(cmd.Args, "rendererType=dom") {
		t.Fatalf("expected dom renderer, got %v", cmd.Args)
	}

	t.Setenv("VHS_RENDERER", "bogus")
	cmd = buildTtyCmd(context.Background(), 1976, shell)
	if !argsContain(cmd.Args, "rendererType=canvas") {
		t.Fatalf("expected canvas fallback for unknown renderer, got %v", cmd.Args)
	}
}
