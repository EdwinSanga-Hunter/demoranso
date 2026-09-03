//go:build !windows
// +build !windows

// Package tor provides a wrapper around the tor proxy command
package tor

import "os/exec"

// Tor wraps the tor command line
type Tor struct {
	RootPath string    // RootPath is the path where a tor bundle would be extracted
	Cmd      *exec.Cmd // Cmd is the tor proxy command
}

// New returns a new Tor instance
func New(rootPath string) *Tor {
	return &Tor{RootPath: rootPath}
}

// DownloadAndExtract is a no-op outside windows.
// The client will instead use a locally running tor service
// (e.g. `sudo service tor start`), reachable at 127.0.0.1:9050.
func (t *Tor) DownloadAndExtract() error {
	return nil
}

// Start is a no-op outside windows
func (t *Tor) Start() error {
	return nil
}

// GetExecutable returns the tor command name
func (t *Tor) GetExecutable() string {
	return "tor"
}

// Kill kill the tor process, if any
func (t *Tor) Kill() error {
	if t.Cmd != nil {
		return t.Cmd.Process.Kill()
	}
	return nil
}

// Clean is a no-op outside windows
func (t *Tor) Clean() error {
	return nil
}
