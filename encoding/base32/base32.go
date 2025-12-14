// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package base32

import (
	"encoding/base32"

	"github.com/demurky/monad/encoding"
)

func Encode[Output, Input ~string | ~[]byte](src Input) Output {
	return encoding.Encode[Output](base32.StdEncoding, src)
}
func EncodeHex[Output, Input ~string | ~[]byte](src Input) Output {
	return encoding.Encode[Output](base32.HexEncoding, src)
}

func Decode[Output, Input ~string | ~[]byte](src Input) (Output, error) {
	return encoding.Decode[Output](base32.StdEncoding, src)
}
func DecodeHex[Output, Input ~string | ~[]byte](src Input) (Output, error) {
	return encoding.Decode[Output](base32.HexEncoding, src)
}
