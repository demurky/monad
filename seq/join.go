// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package seq

import (
	"bytes"
	"iter"
	"strings"
)

func Join(s iter.Seq[string], sep string) string {
	var sb strings.Builder
	first := true
	for v := range s {
		if !first {
			sb.WriteString(sep)
		}
		sb.WriteString(v)
		first = false
	}
	return sb.String()
}

func JoinBytes(s iter.Seq[[]byte], sep []byte) []byte {
	var b bytes.Buffer
	first := true
	for v := range s {
		if !first {
			b.Write(sep)
		}
		b.Write(v)
		first = false
	}
	return b.Bytes()
}
