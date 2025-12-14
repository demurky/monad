// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package noalloc

import "unsafe"

func ConvertBytes[T ~string | ~[]byte](b []byte) T {
	var v T
	if len(b) == 0 {
		return v
	}
	switch any(v).(type) {
	case []byte:
		return T(b)
	case string:
		return T(unsafe.String(&b[0], len(b)))
	default:
		panic("unreachable")
	}
}
