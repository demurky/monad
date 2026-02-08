// Copyright 2025 the toolbox authors.
// SPDX-License-Identifier: Apache-2.0

// Package oslite offers allocation-free implementations
// of some of the [os] package's functionality.
//
// For now, only Linux has custom implementations,
// and we fall back to [os] on other platforms.
//
// This package is ideal for tasks like repeatedly opening files
// in in-memory filesystems (e.g. Linux's procfs).
package oslite

import "os"

type File struct {
	fd   int
	file *os.File
}

func Open(path string) (f File, err error) {
	return f, f.Open(path)
}

func Create(path string) (f File, err error) {
	return f, f.Create(path)
}

func OpenFile(path string, flag int, perm os.FileMode) (f File, err error) {
	return f, f.OpenFile(path, flag, perm)
}
