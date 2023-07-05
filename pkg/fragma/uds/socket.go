// Package uds contains Unix Domain Socket utilities.
package uds

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func IsAvailable() error {
	fpath, err := NewSocket()
	if err != nil {
		return fmt.Errorf("create socket: %v", err)
	}
	defer func() { _ = CleanupSocket(fpath) }()

	lis, err := net.Listen("unix", fpath)
	if err != nil {
		return fmt.Errorf("listen on socket: %v", err)
	}
	defer lis.Close()

	return nil
}

func NewSocket() (string, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("create temporary directory: %v", err)
	}

	fpath := filepath.Join(dir, "auditum.sock")
	return fpath, nil
}

func CleanupSocket(fpath string) error {
	return os.RemoveAll(filepath.Dir(fpath))
}
