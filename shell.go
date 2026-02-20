package main

import (
	"fmt"
	"strings"
)

// Supported shells of VH.
const (
	bash       = "bash"
	cmdexe     = "cmd"
	fish       = "fish"
	nushell    = "nu"
	osh        = "osh"
	powershell = "powershell"
	pwsh       = "pwsh"
	xonsh      = "xonsh"
	zsh        = "zsh"
)

// DefaultPromptColor is the default color for the shell prompt.
const DefaultPromptColor = "#5B56E0"

// DefaultPrompt is the default prompt symbol.
const DefaultPrompt = ">"

// Shell is a type that contains a prompt and the command to set up the shell.
type Shell struct {
	Name string
}

// ShellConfig returns the shell configuration with the given prompt and color.
func ShellConfig(name, promptColor, prompt string) (env []string, command []string) {
	// Parse hex color to RGB components
	r, g, b := hexToRGB(promptColor)
	hexNoHash := strings.TrimPrefix(promptColor, "#")

	switch name {
	case bash:
		return []string{
				fmt.Sprintf("PS1=\\[\\e[38;2;%d;%d;%dm\\]%s \\[\\e[0m\\]", r, g, b, prompt),
				"BASH_SILENCE_DEPRECATION_WARNING=1",
			},
			[]string{"bash", "--noprofile", "--norc", "--login", "+o", "history"}
	case zsh:
		return []string{fmt.Sprintf(`PROMPT=%%F{#%s}%s %%F{reset_color}`, hexNoHash, prompt)},
			[]string{"zsh", "--histnostore", "--no-rcs"}
	case fish:
		return nil, []string{
			"fish",
			"--login",
			"--no-config",
			"--private",
			"-C", "function fish_greeting; end",
			"-C", fmt.Sprintf(`function fish_prompt; set_color %s; echo -n "%s "; set_color normal; end`, hexNoHash, prompt),
		}
	case powershell:
		return nil, []string{
			"powershell",
			"-NoLogo",
			"-NoExit",
			"-NoProfile",
			"-Command",
			fmt.Sprintf(`Set-PSReadLineOption -HistorySaveStyle SaveNothing; function prompt { Write-Host '%s' -NoNewLine -ForegroundColor ([System.Drawing.Color]::FromArgb(%d,%d,%d)); return ' ' }`, prompt, r, g, b),
		}
	case pwsh:
		return nil, []string{
			"pwsh",
			"-Login",
			"-NoLogo",
			"-NoExit",
			"-NoProfile",
			"-Command",
			fmt.Sprintf(`Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor ([System.Drawing.Color]::FromArgb(%d,%d,%d)) -NoNewLine '%s'; return ' ' }`, r, g, b, prompt),
		}
	case cmdexe:
		return nil, []string{"cmd.exe", "/k", fmt.Sprintf("prompt=%s ", prompt)}
	case nushell:
		return nil, []string{"nu", "--execute", fmt.Sprintf("$env.PROMPT_COMMAND = {'\033[;38;2;%d;%d;%dm%s\033[m '}; $env.PROMPT_COMMAND_RIGHT = {''}", r, g, b, prompt)}
	case osh:
		return []string{fmt.Sprintf("PS1=\\[\\e[38;2;%d;%d;%dm\\]%s \\[\\e[0m\\]", r, g, b, prompt)},
			[]string{"osh", "--norc"}
	case xonsh:
		return nil, []string{"xonsh", "--no-rc", "-D", fmt.Sprintf("PROMPT=\033[;38;2;%d;%d;%dm%s\033[m ", r, g, b, prompt)}
	default:
		return nil, nil
	}
}

// hexToRGB converts a hex color string to RGB components.
func hexToRGB(hex string) (r, g, b int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 6 {
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	}
	return
}

// Shells contains a mapping from shell names to their Shell struct.
var Shells = map[string]Shell{
	bash:       {Name: bash},
	zsh:        {Name: zsh},
	fish:       {Name: fish},
	powershell: {Name: powershell},
	pwsh:       {Name: pwsh},
	cmdexe:     {Name: cmdexe},
	nushell:    {Name: nushell},
	osh:        {Name: osh},
	xonsh:      {Name: xonsh},
}
