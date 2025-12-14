// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package moreio

import (
	"io"
	"sync"
	"unicode/utf8"
	"unsafe"
)

var (
	bytePool = sync.Pool{
		New: func() any {
			var buf [1]byte
			return &buf
		},
	}
	runePool = sync.Pool{
		New: func() any {
			var buf [utf8.UTFMax]byte
			return &buf
		},
	}
)

// AdapterWriter is returned by [NewAdapterWriter].
// See it's documentation for details.
type AdapterWriter struct {
	w io.Writer
	s io.StringWriter
	b io.ByteWriter
	r runeWriter

	c   config
	mu  sync.Mutex
	buf [utf8.UTFMax]byte
}

// There is no io.RuneWriter in the stdlib:
// https://github.com/golang/go/issues/71027
type runeWriter interface {
	WriteRune(r rune) (int, error)
}

// NewAdapterWriter returns a writer
// that forwards the WriteString, WriteByte and WriteRune method calls
// if they're implemented by w,
// otherwise it implements them on top of w.Write().
//
// Use the [WithThreadSafe] option
// if you want to call any of the methods concurrently.
func NewAdapterWriter(w io.Writer, opts ...Option) *AdapterWriter {

	if a, ok := w.(*AdapterWriter); ok {
		return a
	}

	a := &AdapterWriter{w: w}
	a.s, _ = w.(io.StringWriter)
	a.b, _ = w.(io.ByteWriter)
	a.r, _ = w.(runeWriter)

	for _, fn := range opts {
		fn(&a.c)
	}

	return a
}

func (a *AdapterWriter) Write(b []byte) (int, error) {
	return a.w.Write(b)
}

func (a *AdapterWriter) WriteString(s string) (int, error) {
	if a.s != nil {
		return a.s.WriteString(s)
	}
	return a.w.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// When WriteByte and WriteRune try to call a.W.Write(),
// they're basically passing their stack-allocated input (`c byte` or `r rune`)
// to an unknown function.

func (a *AdapterWriter) WriteByte(c byte) error {

	if a.b != nil {
		return a.b.WriteByte(c)
	}

	if !a.c.threadSafe {
		a.buf[0] = c
		_, err := a.w.Write(a.buf[:1])
		return err
	}

	buf := bytePool.Get().(*[1]byte)
	defer bytePool.Put(buf)

	buf[0] = c
	_, err := a.w.Write(buf[:1])
	return err
}

func (a *AdapterWriter) WriteRune(r rune) (int, error) {

	if a.r != nil {
		return a.r.WriteRune(r)
	}

	if !a.c.threadSafe {
		n := utf8.EncodeRune(a.buf[:], r)
		return a.w.Write(a.buf[:n])
	}

	buf := runePool.Get().(*[utf8.UTFMax]byte)
	defer runePool.Put(buf)

	n := utf8.EncodeRune(buf[:], r)
	return a.w.Write(buf[:n])
}
