package main

import (
	"strings"
	"testing"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex     string
		r, g, b int
	}{
		{"#5B56E0", 91, 86, 224},
		{"5B56E0", 91, 86, 224},
		{"#F6821F", 246, 130, 31},
		{"#000000", 0, 0, 0},
		{"#FFFFFF", 255, 255, 255},
		{"", 0, 0, 0},
		{"#FFF", 0, 0, 0}, // invalid length
	}

	for _, tc := range tests {
		r, g, b := hexToRGB(tc.hex)
		if r != tc.r || g != tc.g || b != tc.b {
			t.Errorf("hexToRGB(%q) = (%d,%d,%d), want (%d,%d,%d)",
				tc.hex, r, g, b, tc.r, tc.g, tc.b)
		}
	}
}

func TestShellConfigReturnsNonNil(t *testing.T) {
	// Every known shell should return a non-nil command from ShellConfig.
	shellNames := []string{
		bash, zsh, fish, powershell, pwsh, cmdexe, nushell, osh, xonsh,
	}

	for _, name := range shellNames {
		t.Run(name, func(t *testing.T) {
			_, command := ShellConfig(name, DefaultPromptColor, DefaultPrompt)
			if len(command) == 0 {
				t.Errorf("ShellConfig(%q, %q) returned empty command", name, DefaultPromptColor)
			}
		})
	}
}

func TestShellConfigCustomColor(t *testing.T) {
	// Shells that use RGB values should embed the custom colour.
	// bash uses PS1 with ANSI 38;2;R;G;B escape, so changing the colour
	// should produce different RGB values in the env string.
	env, _ := ShellConfig(bash, "#FF8000", DefaultPrompt)
	if len(env) == 0 {
		t.Fatal("expected env for bash")
	}
	// #FF8000 = 255,128,0
	if !strings.Contains(env[0], "255;128;0") {
		t.Errorf("bash PS1 does not contain expected RGB values for #FF8000: %s", env[0])
	}

	// zsh uses hex in the PROMPT string
	env, _ = ShellConfig(zsh, "#FF8000", DefaultPrompt)
	if len(env) == 0 {
		t.Fatal("expected env for zsh")
	}
	if !strings.Contains(env[0], "FF8000") {
		t.Errorf("zsh PROMPT does not contain hex colour FF8000: %s", env[0])
	}
}

func TestShellConfigDefaultColor(t *testing.T) {
	// With the default colour, bash should produce the original RGB values.
	// #5B56E0 = 91,86,224
	env, _ := ShellConfig(bash, DefaultPromptColor, DefaultPrompt)
	if len(env) == 0 {
		t.Fatal("expected env for bash")
	}
	if !strings.Contains(env[0], "91;86;224") {
		t.Errorf("bash PS1 with default colour does not contain expected RGB: %s", env[0])
	}
}

func TestShellConfigCustomPrompt(t *testing.T) {
	// Every shell should embed the custom prompt symbol in its config.
	shellNames := []string{
		bash, zsh, fish, powershell, pwsh, cmdexe, nushell, osh, xonsh,
	}

	for _, name := range shellNames {
		t.Run(name, func(t *testing.T) {
			env, command := ShellConfig(name, DefaultPromptColor, "λ")
			combined := strings.Join(env, " ") + " " + strings.Join(command, " ")
			if !strings.Contains(combined, "λ") {
				t.Errorf("ShellConfig(%q) with prompt λ does not contain the symbol: env=%v command=%v", name, env, command)
			}
		})
	}
}

func TestShellConfigUnknownShell(t *testing.T) {
	env, command := ShellConfig("unknown-shell", DefaultPromptColor, DefaultPrompt)
	if env != nil || command != nil {
		t.Errorf("expected nil for unknown shell, got env=%v command=%v", env, command)
	}
}
