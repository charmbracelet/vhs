package main

import (
	"strings"
	"testing"
)

func TestShellPromptMarker(t *testing.T) {
	// Every shell configuration should embed the OSC 133;A prompt marker
	// (FinalTerm shell integration) so that AwaitPrompt can detect when a
	// command has finished.
	//
	// The marker format varies by shell:
	//   - Most shells: \e]133;A\a (ESC ] 133;A BEL)
	//   - cmd.exe: $E]133;A$E\ (using ST terminator instead of BEL)
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

			if !strings.Contains(combined, "133") {
				t.Errorf("Shell %q does not contain OSC 133;A marker.\nenv: %v\ncommand: %v", name, shell.Env, shell.Command)
			}
		})
	}
}
