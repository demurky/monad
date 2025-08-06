// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package seq

import (
	"bufio"
	"errors"
	"fmt"
	"iter"
	"os"
	"strings"
)

func Join(sep string, seq iter.Seq[string]) string {
	var sb strings.Builder
	first := true
	for s := range seq {
		if !first {
			sb.WriteString(sep)
		}
		sb.WriteString(s)
		first = false
	}
	return sb.String()
}

func Map[T any](f func(T) T, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(v T) bool) {
		for v := range seq {
			if !yield(f(v)) {
				break
			}
		}
	}
}

func FlatMap[T any](f func(T) []T, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(v T) bool) {
	outer:
		for v := range seq {
			for _, r := range f(v) {
				if !yield(r) {
					break outer
				}
			}
		}
	}
}

func FileLines(errPtr *error, paths ...string) iter.Seq[string] {
	*errPtr = nil
	errList := make([]error, 0)
	return func(yield func(line string) bool) {
	outer:
		for _, path := range paths {
			file, err := os.Open(path)
			if err != nil {
				*errPtr = fmt.Errorf("could not open file %q: %w", path, err)
				return
			}
			s := bufio.NewScanner(file)
			for s.Scan() {
				line := s.Text()
				if len(line) == 0 {
					continue
				}
				if !yield(line) {
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
