// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package base64

import (
	"encoding/base64"

	"github.com/demurky/monad/encoding"
)

func Encode[Output, Input ~string | ~[]byte](src Input) Output {
	return encoding.Encode[Output](base64.StdEncoding, src)
}
func EncodeURL[Output, Input ~string | ~[]byte](src Input) Output {
	return encoding.Encode[Output](base64.URLEncoding, src)
}
func EncodeRaw[Output, Input ~string | ~[]byte](src Input) Output {
	return encoding.Encode[Output](base64.RawStdEncoding, src)
}
func EncodeRawURL[Output, Input ~string | ~[]byte](src Input) Output {
	return encoding.Encode[Output](base64.RawURLEncoding, src)
}

func Decode[Output, Input ~string | ~[]byte](src Input) (Output, error) {
	return encoding.Decode[Output](base64.StdEncoding, src)
}
func DecodeURL[Output, Input ~string | ~[]byte](src Input) (Output, error) {
	return encoding.Decode[Output](base64.URLEncoding, src)
}
func DecodeRaw[Output, Input ~string | ~[]byte](src Input) (Output, error) {
	return encoding.Decode[Output](base64.RawStdEncoding, src)
}
func DecodeRawURL[Output, Input ~string | ~[]byte](src Input) (Output, error) {
	return encoding.Decode[Output](base64.RawURLEncoding, src)
}
