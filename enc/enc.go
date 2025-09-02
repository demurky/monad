// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package enc

import (
	"encoding/base64"
)

type Encoding interface {
	EncodedLen(n int) int
	DecodedLen(n int) int
	Encode(dst, src []byte)
	Decode(dst, src []byte) (int, error)
}

func Encode[Dst, Src ~[]byte | ~string](enc Encoding, src Src) Dst {
	dst := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(dst, []byte(src))
	return Dst(dst)
}

func Decode[Dst, Src ~[]byte | ~string](enc Encoding, src Src) (Dst, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	_, err := enc.Decode(dst, []byte(src))
	return Dst(dst), err
}
