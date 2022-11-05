//go:build windows
// +build windows

package main

import "os/exec"

var Shells = map[string]LazyShell{
	cmdexe: func() Shell {
		return Shell{
			EntryPoint: "cmd",
			Prompt:     "$g",
			Command:    ` cls && set prompt=%s && cls`,
		}
	},
	pwsh: func() Shell {
		if _, err := exec.LookPath("pwsh"); err == nil {
			return Shell{
				EntryPoint: "powershell",
				Prompt:     "Function prompt {Write-Host \"> \" -ForegroundColor Blue -NoNewLine; return \"`0\" }",
				Command:    ` clear; pwsh -Login -NoLogo -NoExit -Command 'Set-PSReadLineOption -HistorySaveStyle SaveNothing; %s'`,
			}
		}

		return Shell{
			EntryPoint: "powershell",
			Prompt:     "Function prompt {Write-Host \"> \" -ForegroundColor Blue -NoNewLine; return \"`0\" }",
			Command:    ` clear; powershell -NoLogo -NoExit -Command 'Set-PSReadLineOption -HistorySaveStyle SaveNothing; %s'`,
		}
	},
}

var defaultShell = cmdexe
