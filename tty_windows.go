//go:build windows
// +build windows

package main

import "golang.org/x/sys/windows"

func defaultShellWithArgs() []string {
	shell := "cmd"
	major, _, _ := windows.RtlGetNtVersionNumbers()

	if major >= 10 {
		shell = "powershell"
	}

	return []string{
		shell,
	}
}
