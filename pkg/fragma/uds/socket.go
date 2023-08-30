// Copyright 2023 Igor Zibarev
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
