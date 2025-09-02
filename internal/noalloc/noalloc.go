// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package noalloc

import "unsafe"

func FromBytes[Dst ~[]byte | ~string](b []byte) Dst {
	var v Dst
	if len(b) == 0 {
		return v
	}
	switch any(v).(type) {
	case []byte:
		return Dst(b)
	case string:
		return Dst(unsafe.String(&b[0], len(b)))
	default:
		panic("unreachable")
	}
}
