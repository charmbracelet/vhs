package main

import (
	"fmt"
	"strconv"
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

// hexColorLen is the number of characters in a hex color, sans the leading "#".
const hexColorLen = 6

// Shell is a type that contains a prompt and the command to set up the shell.
type Shell struct {
	Name string
}

// ShellConfig returns the shell configuration with the given prompt color.
func ShellConfig(name, promptColor string) (env []string, command []string) {
	// Parse the hex color once so every shell gets the same value. Fall back
	// to the default color if it can't be parsed, so bad input never reaches
	// a shell.
	r, g, b, err := hexToRGB(promptColor)
	if err != nil {
		r, g, b, _ = hexToRGB(DefaultPromptColor)
	}
	hexNoHash := fmt.Sprintf("%02X%02X%02X", r, g, b)

	switch name {
	case bash:
		return []string{
				fmt.Sprintf("PS1=\\[\\e[38;2;%d;%d;%dm\\]> \\[\\e[0m\\]", r, g, b),
				"BASH_SILENCE_DEPRECATION_WARNING=1",
			},
			[]string{"bash", "--noprofile", "--norc", "--login", "+o", "history"}
	case zsh:
		return []string{fmt.Sprintf(`PROMPT=%%F{#%s}> %%F{reset_color}`, hexNoHash)},
			[]string{"zsh", "--histnostore", "--no-rcs"}
	case fish:
		return nil, []string{
			"fish",
			"--login",
			"--no-config",
			"--private",
			"-C", "function fish_greeting; end",
			"-C", fmt.Sprintf(`function fish_prompt; set_color %s; echo -n "> "; set_color normal; end`, hexNoHash),
		}
	case powershell:
		return nil, []string{
			"powershell",
			"-NoLogo",
			"-NoExit",
			"-NoProfile",
			"-Command",
			fmt.Sprintf(`Set-PSReadLineOption -HistorySaveStyle SaveNothing; function prompt { Write-Host '>' -NoNewLine -ForegroundColor ([System.Drawing.Color]::FromArgb(%d,%d,%d)); return ' ' }`, r, g, b),
		}
	case pwsh:
		return nil, []string{
			"pwsh",
			"-Login",
			"-NoLogo",
			"-NoExit",
			"-NoProfile",
			"-Command",
			fmt.Sprintf(`Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor ([System.Drawing.Color]::FromArgb(%d,%d,%d)) -NoNewLine '>'; return ' ' }`, r, g, b),
		}
	case cmdexe:
		return nil, []string{"cmd.exe", "/k", "prompt=^> "}
	case nushell:
		return nil, []string{"nu", "--execute", fmt.Sprintf("$env.PROMPT_COMMAND = {'\033[;38;2;%d;%d;%dm>\033[m '}; $env.PROMPT_COMMAND_RIGHT = {''}", r, g, b)}
	case osh:
		return []string{fmt.Sprintf("PS1=\\[\\e[38;2;%d;%d;%dm\\]> \\[\\e[0m\\]", r, g, b)},
			[]string{"osh", "--norc"}
	case xonsh:
		return nil, []string{"xonsh", "--no-rc", "-D", fmt.Sprintf("PROMPT=\033[;38;2;%d;%d;%dm>\033[m ", r, g, b)}
	default:
		return nil, nil
	}
}

// hexToRGB converts a hex color string, with or without a leading "#", to RGB
// components. It returns an error if the string is not a 6-digit hex color.
func hexToRGB(hex string) (r, g, b int, err error) {
	digits := strings.TrimPrefix(hex, "#")
	n, parseErr := strconv.ParseUint(digits, 16, 32)
	if len(digits) != hexColorLen || parseErr != nil {
		return 0, 0, 0, fmt.Errorf("%q is not a valid hex color, expected 6 hex digits like \"#5B56E0\"", hex)
	}
	return int(n >> 16 & 0xFF), int(n >> 8 & 0xFF), int(n & 0xFF), nil
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
