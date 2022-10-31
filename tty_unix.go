//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

func defaultShellWithArgs() []string {
	return []string{
		"bash", "--login",
	}
}
