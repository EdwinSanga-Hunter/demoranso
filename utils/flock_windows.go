//go:build windows
// +build windows

package utils

// FileIsLocked is a no-op on windows
func FileIsLocked(path string) bool {
	return false
}
