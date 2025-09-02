// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package text

import "unsafe"

func Wrap[Dst, Src ~[]byte | ~string](src Src, limit int) Dst {
	dst := make([]byte, 0, len(src)+(len(src)/limit))
	if limit <= 0 {
		panic("Wrap: limit <= 0")
	}
	n := limit
	i := 0
	for len(src) > 0 {
		dst = append(dst, src[i])
		i++
		n--
		if n == 0 {
			dst = append(dst, '\n')
			n = limit
		}
	}
	return convertBytes[Dst](dst)
}

func convertBytes[T ~[]byte | ~string](b []byte) T {
	var t T
	if len(b) == 0 {
		return t
	}
	switch any(t).(type) {
	case []byte:
		return T(b)
	case string:
		return T(unsafe.String(&b[0], len(b)))
	default:
		panic("unreachable")
	}
}
