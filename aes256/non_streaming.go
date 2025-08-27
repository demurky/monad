// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package aes256

import (
	"bytes"
)

func Encrypt(plaintext, password []byte) []byte {
	ciphertext := bytes.NewBuffer(make([]byte, 0, len(plaintext)+100))
	w := NewWriter(ciphertext, password)
	w.Write(plaintext)
	w.Close()
	return ciphertext.Bytes()
}

func Decrypt(ciphertext, password []byte) ([]byte, error) {
	r := NewReader(bytes.NewReader(ciphertext), password)
	plaintext := make([]byte, len(ciphertext))
	n, err := r.Read(plaintext)
	return plaintext[:n], err
}
