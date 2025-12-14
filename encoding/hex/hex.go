// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package hex

import (
	"encoding/hex"

	"github.com/demurky/monad/encoding"
)

var Encoding encoding.Encoding = e{}

type e struct{}

func (e) EncodedLen(n int) int {
	return hex.EncodedLen(n)
}

func (e) DecodedLen(n int) int {
	return hex.DecodedLen(n)
}

func (e) Encode(dst, src []byte) {
	hex.Encode(dst, src)
}

func (e) Decode(dst, src []byte) (int, error) {
	return hex.Decode(dst, src)
}

func Encode[Output, Input ~string | ~[]byte](src Input) Output {
	return encoding.Encode[Output](Encoding, src)
}

func Decode[Output, Input ~string | ~[]byte](src Input) (Output, error) {
	return encoding.Decode[Output](Encoding, src)
}
