// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package moreio

import (
	"io"
)

// ErrorCapturingWriter is returned by [NewErrorCapturingWriter].
// See it's documentation for details.
type ErrorCapturingWriter struct {
	W   *AdapterWriter
	Err error
}

// NewErrorCapturingWriter wraps w so that the first write error encountered
// is stored in [ErrorCapturingWriter.Err] and subsequent writes are no-ops.
//
// It's useful for when many small writes are performed
// and handling the error on each write is overkill.
//
// [ErrorCapturingWriter] provides the functionality of [AdapterWriter]
// as well.
func NewErrorCapturingWriter(w io.Writer, opts ...Option) *ErrorCapturingWriter {
	return &ErrorCapturingWriter{
		W: NewAdapterWriter(w),
	}
}

func (w *ErrorCapturingWriter) Write(b []byte) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	n, err := w.W.Write(b)
	if err != nil {
		w.Err = err
	}
	return n, err
}

func (w *ErrorCapturingWriter) WriteString(s string) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	n, err := w.W.WriteString(s)
	if err != nil {
		w.Err = err
	}
	return n, err
}

func (w *ErrorCapturingWriter) WriteByte(c byte) error {
	if w.Err != nil {
		return w.Err
	}
	err := w.W.WriteByte(c)
	if err != nil {
		w.Err = err
	}
	return err
}

func (w *ErrorCapturingWriter) WriteRune(r rune) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	n, err := w.W.WriteRune(r)
	if err != nil {
		w.Err = err
	}
	return n, err
}
