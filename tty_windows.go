//go:build windows
// +build windows

package main

import "golang.org/x/sys/windows"

const defaultShell = pwsh

func defaultShellWithArgs() []string {
	major, _, _ := windows.RtlGetNtVersionNumbers()
	if major >= 10 {
		return []string{"powershell"}
	}

	return []string{"cmd"}
}
