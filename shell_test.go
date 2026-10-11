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
		{"#5b56e0", 91, 86, 224},
		{"#F6821F", 246, 130, 31},
		{"#000000", 0, 0, 0},
		{"#FFFFFF", 255, 255, 255},
	}

	for _, tc := range tests {
		r, g, b, err := hexToRGB(tc.hex)
		if err != nil {
			t.Errorf("hexToRGB(%q) returned unexpected error: %v", tc.hex, err)
			continue
		}
		if r != tc.r || g != tc.g || b != tc.b {
			t.Errorf("hexToRGB(%q) = (%d,%d,%d), want (%d,%d,%d)",
				tc.hex, r, g, b, tc.r, tc.g, tc.b)
		}
	}
}

func TestHexToRGBInvalid(t *testing.T) {
	invalid := []string{
		"",
		"#",
		"#FFF",      // too short
		"#5B56E",    // missing a digit
		"#5B56E00",  // too long
		"#GGGGGG",   // not hex
		"#12345Z",   // last digit not hex
		"#+12345",   // sign
		"red",       // named colour
		"##5B56E0",  // double hash
		" #5B56E0 ", // surrounding spaces
	}

	for _, hex := range invalid {
		if _, _, _, err := hexToRGB(hex); err == nil {
			t.Errorf("hexToRGB(%q) expected an error, got nil", hex)
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
			_, command := ShellConfig(name, DefaultPromptColor)
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
	env, _ := ShellConfig(bash, "#FF8000")
	if len(env) == 0 {
		t.Fatal("expected env for bash")
	}
	// #FF8000 = 255,128,0
	if !strings.Contains(env[0], "255;128;0") {
		t.Errorf("bash PS1 does not contain expected RGB values for #FF8000: %s", env[0])
	}

	// zsh uses hex in the PROMPT string
	env, _ = ShellConfig(zsh, "#FF8000")
	if len(env) == 0 {
		t.Fatal("expected env for zsh")
	}
	if !strings.Contains(env[0], "FF8000") {
		t.Errorf("zsh PROMPT does not contain hex colour FF8000: %s", env[0])
	}
}

func TestShellConfigPowerShellColor(t *testing.T) {
	// PowerShell should set the colour with an ANSI escape sequence, since
	// Write-Host -ForegroundColor only accepts named console colours.
	_, command := ShellConfig(powershell, "#FF8000")
	prompt := command[len(command)-1]
	if !strings.Contains(prompt, "[38;2;255;128;0m") {
		t.Errorf("powershell prompt does not contain expected ANSI colour for #FF8000: %s", prompt)
	}
	if strings.Contains(prompt, "System.Drawing") {
		t.Errorf("powershell prompt should not use System.Drawing: %s", prompt)
	}
}

func TestShellConfigDefaultColor(t *testing.T) {
	// With the default colour, bash should produce the original RGB values.
	// #5B56E0 = 91,86,224
	env, _ := ShellConfig(bash, DefaultPromptColor)
	if len(env) == 0 {
		t.Fatal("expected env for bash")
	}
	if !strings.Contains(env[0], "91;86;224") {
		t.Errorf("bash PS1 with default colour does not contain expected RGB: %s", env[0])
	}
}

func TestShellConfigInvalidColorFallsBack(t *testing.T) {
	// An invalid colour should give every shell the default colour, rather
	// than passing the bad value through to shells that use hex directly.
	env, _ := ShellConfig(bash, "not-a-color")
	if !strings.Contains(env[0], "91;86;224") {
		t.Errorf("bash PS1 should use the default colour for invalid input: %s", env[0])
	}

	env, _ = ShellConfig(zsh, "not-a-color")
	if !strings.Contains(env[0], "%F{#5B56E0}") {
		t.Errorf("zsh PROMPT should use the default colour for invalid input: %s", env[0])
	}

	_, command := ShellConfig(fish, "not-a-color")
	if !strings.Contains(strings.Join(command, " "), "set_color 5B56E0") {
		t.Errorf("fish prompt should use the default colour for invalid input: %v", command)
	}
}

func TestShellConfigUnknownShell(t *testing.T) {
	env, command := ShellConfig("unknown-shell", DefaultPromptColor)
	if env != nil || command != nil {
		t.Errorf("expected nil for unknown shell, got env=%v command=%v", env, command)
	}
}
