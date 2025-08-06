// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package seq

import "iter"

func Values[T any](values ...T) iter.Seq[T] {
	return func(yield func(v T) bool) {
		for _, v := range values {
			if !yield(v) {
				break
			}
		}
	}
}
