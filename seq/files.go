// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package seq

import (
	"bufio"
	"errors"
	"fmt"
	"iter"
	"os"
)

func FileLines(errPtr *error, paths ...string) iter.Seq[string] {
	return ScanFiles(errPtr, bufio.ScanLines, paths...)
}

func ScanFiles(errPtr *error, split bufio.SplitFunc, paths ...string) iter.Seq[string] {
	*errPtr = nil
	errList := make([]error, 0)
	return func(yield func(token string) bool) {
	outer:
		for _, path := range paths {
			file, err := os.Open(path)
			if err != nil {
				*errPtr = fmt.Errorf("could not open file %q: %w", path, err)
				return
			}
			s := bufio.NewScanner(file)
			s.Split(split)
			for s.Scan() {
				token := s.Text()
				if len(token) == 0 {
					continue
				}
				if !yield(token) {
					file.Close()
					break outer
				}
			}
			err = s.Err()
			if err != nil {
				errList = append(errList, fmt.Errorf(
					"could not scan file %q: %w",
					path, err,
				))
			}
			file.Close()
		}
		*errPtr = errors.Join(errList...)
	}
}
