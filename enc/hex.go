// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package enc

import "encoding/hex"

var HexEncoding Encoding = hexEncoding{}

type hexEncoding struct{}

func (hexEncoding) EncodedLen(n int) int                { return hex.EncodedLen(n) }
func (hexEncoding) DecodedLen(n int) int                { return hex.DecodedLen(n) }
func (hexEncoding) Encode(dst, src []byte)              { hex.Encode(dst, src) }
func (hexEncoding) Decode(dst, src []byte) (int, error) { return hex.Decode(dst, src) }

func EncHex[Dst, Src ~[]byte | ~string](src Src) Dst {
	return Encode[Dst](HexEncoding, src)
}

func DecHex[Dst, Src ~[]byte | ~string](src Src) (Dst, error) {
	return Decode[Dst](HexEncoding, src)
}
