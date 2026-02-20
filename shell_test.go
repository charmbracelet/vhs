package main

import (
	"strings"
	"testing"
)

func TestShellPromptMarker(t *testing.T) {
	// Every shell configuration should embed the OSC 7777 prompt marker
	// so that AwaitPrompt can detect when a command has finished.
	//
	// The marker format varies by shell:
	//   - Most shells: \e]7777;\a (ESC ] 7777 ; BEL)
	//   - cmd.exe: $E]7777;$E\ (using ST terminator instead of BEL)
	shellNames := []string{
		bash,
		zsh,
		fish,
		powershell,
		pwsh,
		cmdexe,
		nushell,
		osh,
		xonsh,
	}

	for _, name := range shellNames {
		t.Run(name, func(t *testing.T) {
			shell, ok := Shells[name]
			if !ok {
				t.Fatalf("Shell %q not found in Shells map", name)
			}

			// Combine env and command into a single string to search
			combined := strings.Join(shell.Env, " ") + " " + strings.Join(shell.Command, " ")

			if !strings.Contains(combined, "7777") {
				t.Errorf("Shell %q does not contain OSC 7777 marker.\nenv: %v\ncommand: %v", name, shell.Env, shell.Command)
			}
		})
	}
}
