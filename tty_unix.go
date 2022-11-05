//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

func defaultEntryPoint() []string {
	return []string{"bash", "--login"}
}
