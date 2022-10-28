//go:build windows
// +build windows

package main

// Windows doesn't support UID and GID, so we need to skip this.
func dropUserPrivileges(gid int, uid int) error {
	return nil
}
