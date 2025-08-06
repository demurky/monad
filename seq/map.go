// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package seq

import (
	"iter"
)

func Map[T any](s iter.Seq[T], f func(T) T) iter.Seq[T] {
	return func(yield func(v T) bool) {
		for v := range s {
			if !yield(f(v)) {
				break
			}
		}
	}
}

func FlatMap[T any](s iter.Seq[T], f func(T) []T) iter.Seq[T] {
	return func(yield func(v T) bool) {
	outer:
		for v := range s {
			for _, r := range f(v) {
				if !yield(r) {
					break outer
				}
			}
		}
	}
}
