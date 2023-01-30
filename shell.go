package main

import (
	_ "embed"
	"os"
)

//go:embed rc.bash
var bashrc []byte

// Supported shells of VHS
const (
	bash       = "bash"
	cmdexe     = "cmd"
	fish       = "fish"
	powershell = "powershell"
	pwsh       = "pwsh"
	zsh        = "zsh"
)

// Shell is a type that contains a prompt and the command to set up the shell.
type Shell struct {
	Command func() ([]string, error)
	Env     []string
}

func writeRc(name string, content []byte) (string, error) {
	f, err := os.CreateTemp("", "vhs-"+name)
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// Shells contains a mapping from shell names to their Shell struct.
var Shells = map[string]Shell{
	bash: {
		Command: func() ([]string, error) {
			path, err := writeRc("bashrc", bashrc)
			if err != nil {
				return nil, err
			}
			return []string{"bash", "--noprofile", "--rcfile", path}, nil
		},
	},
	zsh: {
		Env: []string{`PROMPT=%F{#5B56E0}> %F{reset_color}`},
		Command: func() ([]string, error) {
			return []string{"zsh", "--histnostore", "--no-rcs"}, nil
		},
	},
	fish: {
		Command: func() ([]string, error) {
			return []string{
				"fish",
				"--login",
				"--no-config",
				"--private",
				"-C", "function fish_greeting; end",
				"-C", `function fish_prompt; echo -e "$(set_color 5B56E0)> $(set_color normal)"; end`,
			}, nil
		},
	},
	powershell: {
		Command: func() ([]string, error) {
			const cmd = `Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor Blue -NoNewLine '>'; return ' ' }`
			return []string{
				"powershell",
				"-Login",
				"-NoLogo",
				"-NoExit",
				"-NoProfile",
				"-Command", cmd,
			}, nil
		},
	},
	pwsh: {
		Command: func() ([]string, error) {
			const cmd = `Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor Blue -NoNewLine '>'; return ' ' }`
			return []string{
				"pwsh",
				"-Login",
				"-NoLogo",
				"-NoExit",
				"-NoProfile",
				"-Command", cmd,
			}, nil
		},
	},
	// TODO: oh boy
	// cmdexe: {
	// 	Prompt:  "$g",
	// 	Command: ` cls && set prompt=%s && cls`,
	// },
}
