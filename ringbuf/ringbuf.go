// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package ringbuf

import (
	"slices"
	"unicode/utf8"
)

type Buffer struct {
	buf      []byte
	maxLen   int
	writePos int
}

func New(maxLen int) *Buffer {
	return &Buffer{
		maxLen: maxLen,
	}
}

func (x *Buffer) Write(b []byte) (int, error) {

	if len(b) >= x.maxLen {
		x.writePos = 0
		x.buf = slices.Grow(x.buf, x.maxLen-len(x.buf))[:x.maxLen]
		copy(x.buf, b[len(b)-x.maxLen:])
		return len(b), nil
	}

	rem := x.maxLen - len(x.buf)
	if n := min(len(b), rem); n > 0 {
		x.buf = append(x.buf, b[:n]...)
		b = b[n:]
	}

	for len(b) > 0 {
		n := copy(x.buf[x.writePos:], b)
		b = b[n:]
		x.writePos += n
		if x.writePos == x.maxLen {
			x.writePos = 0
		}
	}

	return len(b), nil
}

func (x *Buffer) WriteString(s string) (int, error) {

	if len(s) >= x.maxLen {
		x.writePos = 0
		x.buf = slices.Grow(x.buf, x.maxLen-len(x.buf))[:x.maxLen]
		copy(x.buf, s[len(s)-x.maxLen:])
		return len(s), nil
	}

	rem := x.maxLen - len(x.buf)
	if n := min(len(s), rem); n > 0 {
		x.buf = append(x.buf, s[:n]...)
		s = s[n:]
	}

	for len(s) > 0 {
		n := copy(x.buf[x.writePos:], s)
		s = s[n:]
		x.writePos += n
		if x.writePos == x.maxLen {
			x.writePos = 0
		}
	}

	return len(s), nil
}

func (x *Buffer) WriteByte(c byte) error {
	x.Write([]byte{c})
	return nil
}

func (x *Buffer) WriteRune(r rune) (int, error) {
	return x.Write(utf8.AppendRune([]byte{}, r))
}

// Bytes returns the contents of the buffer.
//
// The returned slice should only be used for reading,
// since it may alias the buffer content
// at least until the next buffer modification.
//
// Use [Buffer.CloneBytes] if you intend to
// modify the returned slice.
func (x *Buffer) Bytes() []byte {
	if x.writePos == 0 {
		return x.buf
	}
	ret := make([]byte, x.maxLen)
	n := copy(ret, x.buf[x.writePos:])
	copy(ret[n:], x.buf[:x.writePos])
	return ret
}

// CloneBytes is similar to [Buffer.Bytes],
// but the returned slice is a copy of the underlying data.
func (x *Buffer) CloneBytes() []byte {
	if x.writePos == 0 {
		return slices.Clone(x.buf)
	}
	return slices.Concat(x.buf[:x.writePos], x.buf[x.writePos:])
}

// String returns the contents of the buffer as a string.
// If the [Buffer] is a nil pointer,
// it returns "<nil>".
func (x *Buffer) String() string {
	if x == nil {
		return "<nil>"
	}
	return string(x.Bytes())
}

func (x *Buffer) Len() int {
	if x.writePos == 0 {
		return len(x.buf)
	}
	return x.maxLen
}

func (x *Buffer) MaxLen() int {
	return x.maxLen
}

// Truncate grows or shrinks the maximum length of the buffer to n.
func (x *Buffer) Truncate(n int) {
	if x.writePos == 0 {
		x.maxLen = n
		if n < len(x.buf) {
			x.buf = x.buf[len(x.buf)-n:]
		}
		return
	}
	newB := New(n)
	newB.Write(x.Bytes())
	*x = *newB
}
