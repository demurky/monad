// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package enc

import "encoding/base64"

func EncBase64[Dst, Src ~[]byte | ~string](src Src) Dst {
	return Encode[Dst](base64.StdEncoding, src)
}
func EncBase64URL[Dst, Src ~[]byte | ~string](src Src) Dst {
	return Encode[Dst](base64.URLEncoding, src)
}
func EncBase64Raw[Dst, Src ~[]byte | ~string](src Src) Dst {
	return Encode[Dst](base64.RawStdEncoding, src)
}
func EncBase64RawURL[Dst, Src ~[]byte | ~string](src Src) Dst {
	return Encode[Dst](base64.RawURLEncoding, src)
}

func DecBase64[Dst, Src ~[]byte | ~string](src Src) (Dst, error) {
	return Decode[Dst](base64.StdEncoding, src)
}
func DecBase64URL[Dst, Src ~[]byte | ~string](src Src) (Dst, error) {
	return Decode[Dst](base64.URLEncoding, src)
}
func DecBase64Raw[Dst, Src ~[]byte | ~string](src Src) (Dst, error) {
	return Decode[Dst](base64.RawStdEncoding, src)
}
func DecBase64RawURL[Dst, Src ~[]byte | ~string](src Src) (Dst, error) {
	return Decode[Dst](base64.RawURLEncoding, src)
}
