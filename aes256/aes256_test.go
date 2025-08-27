// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package aes256_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koonix/x/aes256"
	"github.com/koonix/x/must"
)

func TestEncryptDecrypt(t *testing.T) {

	input := []byte("hello world")
	password := []byte("mypass123")

	ciphertext := aes256.Encrypt(input, password)
	output, err := aes256.Decrypt(ciphertext, password)
	if err != nil {
		t.Fatalf("could not decrypt: %s", err)
	}

	if diff := cmp.Diff(input, output); diff != "" {
		t.Errorf("incorrect result (-want +got):\n%s", diff)
	}
}

func BenchmarkWrite(b *testing.B) {

	w := aes256.NewWriter(io.Discard, []byte("mypass123"))
	txt := []byte("hello world")

	b.ResetTimer()

	for b.Loop() {
		w.Write(txt)
	}
}

func BenchmarkRead(b *testing.B) {

	password := []byte("mypass123")

	buf := new(bytes.Buffer)
	w := aes256.NewWriter(buf, password)
	w.Write([]byte("hello world"))
	x := &reader{
		b: buf.Bytes(),
	}

	r := aes256.NewReader(x, password)
	buf2 := make([]byte, buf.Len())

	b.ResetTimer()

	for b.Loop() {
		must.Get(r.Read(buf2))
	}
}

type reader struct {
	b []byte
}

func (r *reader) Read(b []byte) (int, error) {
	return copy(b, r.b), nil
}
