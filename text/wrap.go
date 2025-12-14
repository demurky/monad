// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package text

import (
	"github.com/demurky/monad/internal/noalloc"
)

// Wrap hard-wraps the given text.
func Wrap[Output, Input ~string | ~[]byte](src Input, limit int) Output {
	dst := make([]byte, 0, len(src)+(len(src)/limit))
	if limit <= 0 {
		panic("text.Wrap: limit <= 0")
	}
	n := limit
	i := 0
	for i < len(src) {
		dst = append(dst, src[i])
		i++
		n--
		if n == 0 {
			dst = append(dst, '\n')
			n = limit
		}
	}
	return noalloc.ConvertBytes[Output](dst)
}
