// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package encoding

import (
	"encoding/base64"

	"github.com/demurky/monad/internal/noalloc"
)

type Encoding interface {
	EncodedLen(n int) int
	DecodedLen(n int) int
	Encode(dst, src []byte)
	Decode(dst, src []byte) (int, error)
}

func Encode[Output, Input ~string | ~[]byte](enc Encoding, src Input) Output {
	dst := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(dst, []byte(src))
	return noalloc.ConvertBytes[Output](dst)
}

func Decode[Output, Input ~string | ~[]byte](enc Encoding, src Input) (Output, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := enc.Decode(dst, []byte(src))
	return noalloc.ConvertBytes[Output](dst[:n]), err
}
