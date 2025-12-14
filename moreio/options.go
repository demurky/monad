// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package moreio

type Option func(*config)

type config struct {
	threadSafe bool
}

// WithThreadSafe is applicable to [NewAdapterWriter] and [NewErrorCapturingWriter].
func WithThreadSafe() Option {
	return func(c *config) {
		c.threadSafe = true
	}
}
