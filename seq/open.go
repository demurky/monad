// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package seq

import (
	"fmt"
	"iter"
	"os"
)

func Open(paths iter.Seq[string], errPtr *error) iter.Seq[*os.File] {
	return OpenFile(paths, os.O_RDONLY, 0, errPtr)
}

func OpenFile(
	paths iter.Seq[string],
	flag int,
	perm os.FileMode,
	errPtr *error,
) iter.Seq[*os.File] {

	cont := true
	return func(yield func(file *os.File) bool) {
		for path := range paths {
			file, err := os.OpenFile(path, flag, perm)
			if err != nil {
				*errPtr = fmt.Errorf("could not open file %q: %w", path, err)
				return
			}
			cont = yield(file)
			err = file.Close()
			if err != nil {
				*errPtr = fmt.Errorf("could not close file %q: %w", path, err)
				return
			}
			if !cont {
				return
			}
		}
	}
}
