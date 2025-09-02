// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package enc

import "encoding/base32"

func EncBase32[Dst, Src ~[]byte | ~string](src Src) Dst {
	return Encode[Dst](base32.StdEncoding, src)
}
func EncBase32Hex[Dst, Src ~[]byte | ~string](src Src) Dst {
	return Encode[Dst](base32.HexEncoding, src)
}

func DecBase32[Dst, Src ~[]byte | ~string](src Src) (Dst, error) {
	return Decode[Dst](base32.StdEncoding, src)
}
func DecBase32Hex[Dst, Src ~[]byte | ~string](src Src) (Dst, error) {
	return Decode[Dst](base32.HexEncoding, src)
}
