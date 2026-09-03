//go:build !windows
// +build !windows

package utils

import "syscall"

// FileIsLocked reports whether the file at path is currently held
// with an advisory lock by another process (e.g. the server's bolt
// database). Locked files are skipped by the malware so the demo
// can't destroy its own key store.
func FileIsLocked(path string) bool {
	f, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer syscall.Close(f)

	if err := syscall.Flock(f, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return true
	}
	syscall.Flock(f, syscall.LOCK_UN)
	return false
}
