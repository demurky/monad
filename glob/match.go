// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package glob

import (
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func Match(target string, patterns ...string) (bool, error) {
	return match(target, patterns, func(p string, t string) (bool, error) {
		return doublestar.Match(p, t)
	})
}

func match(
	target string,
	patterns []string,
	matchFunc func(pattern string, target string) (bool, error),
) (
	match bool,
	err error,
) {
	for _, ptrn := range patterns {
		positive := true
		if strings.HasPrefix(ptrn, "!") {
			positive = false
			ptrn = ptrn[1:]
		}
		m, err := matchFunc(ptrn, target)
		if err != nil {
			return false, err
		}
		if m {
			match = positive
		}
	}
	return
}
