// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package glob

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type GlobOption func(*globConfig)

type globConfig struct{}

func GlobFS(
	fsys fs.FS,
	patterns []string,
	options ...GlobOption,
) (
	matches []string,
	err error,
) {
	return glob(patterns, func(p string) ([]string, error) {
		return doublestar.Glob(fsys, p)
	})
}

func glob(
	patterns []string,
	globFunc func(pattern string) (matches []string, err error),
) (
	matches []string,
	err error,
) {
	mm := make(map[string]struct{})
	for i, ptrn := range patterns {
		positive := true
		if strings.HasPrefix(ptrn, "!") {
			positive = false
			ptrn = ptrn[1:]
		}
		matches, err := globFunc(ptrn)
		if err != nil {
			return nil, fmt.Errorf("pattern #%d (%q) failed: %w", i, ptrn, err)
		}
		if positive {
			for _, m := range matches {
				mm[m] = struct{}{}
			}
		} else {
			for _, m := range matches {
				delete(mm, m)
			}
		}
	}
	matches = make([]string, 0, len(mm))
	for m := range mm {
		matches = append(matches, m)
	}
	return
}
