//go:build windows
// +build windows

package main

import (
	"os/exec"

	"golang.org/x/sys/windows"
)

var defaultShell = cmdexe

func defaultShellWithArgs() []string {
	major, _, _ := windows.RtlGetNtVersionNumbers()
	if major >= 10 {
		if _, err := exec.LookPath("pwsh"); err == nil {
			defaultShell = pwsh
		} else {
			defaultShell = powershell
		}
	}

	return []string{defaultShell}
}
