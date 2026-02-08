// Copyright 2025 the toolbox authors.
// SPDX-License-Identifier: Apache-2.0

package streamcrypt

import (
	"bytes"
	"io"
	"slices"
)

func Encrypt(
	plaintext []byte,
	password []byte,
	options ...Option,
) []byte {
	ciphertext := bytes.NewBuffer(make([]byte, 0, len(plaintext)+100))
	w := NewEncryptor(ciphertext, password, options...)
	w.Write(plaintext)
	w.Close()
	return slices.Clip(ciphertext.Bytes())
}

func Decrypt(ciphertext []byte, passFunc PasswordFunc, options ...Option) ([]byte, error) {
	r := NewDecryptor(bytes.NewReader(ciphertext), passFunc, options...)
	plaintext := make([]byte, len(ciphertext))
	n, err := r.Read(plaintext)
	if err != io.EOF {
		if err != nil {
			return nil, err
		}
		err = r.Close()
	}
	return slices.Clip(plaintext[:n]), err
}
