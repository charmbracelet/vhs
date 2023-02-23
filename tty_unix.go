//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

const defaultShell = bash

func defaultShellWithArgs() []string {
	return []string{
		"sh",
	}
}
