// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package ringbuf_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/demurky/monad/container/ringbuf"
)

type tests struct {
	bufferSize int
	tests      []test
}

type test struct {
	line  int
	input []byte
	bytes []byte
	len   int
}

var testCases = []tests{
	{
		0,
		[]test{
			{
				line(),
				[]byte{},
				nil,
				0,
			},
			{
				line(),
				[]byte{1, 2, 3},
				nil,
				0,
			},
		},
	},
	{
		5,
		[]test{
			{
				line(),
				[]byte{},
				nil,
				0,
			},
			{
				line(),
				[]byte{1},
				[]byte{1},
				1,
			},
			{
				line(),
				[]byte{2, 3},
				[]byte{1, 2, 3},
				3,
			},
			{
				line(),
				[]byte{4, 5, 6},
				[]byte{2, 3, 4, 5, 6},
				5,
			},
			{
				line(),
				[]byte{7, 8, 9},
				[]byte{5, 6, 7, 8, 9},
				5,
			},
			{
				line(),
				[]byte{},
				[]byte{5, 6, 7, 8, 9},
				5,
			},
			{
				line(),
				[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				[]byte{8, 9, 10, 11, 12},
				5,
			},
		},
	},
}

func TestRingBuf(t *testing.T) {
	for i, test := range testCases {
		b1 := ringbuf.New(test.bufferSize)
		b2 := ringbuf.New(test.bufferSize)

		for ii, test := range test.tests {
			t.Run(fmt.Sprintf("test%d.%d-line%d", i, ii, test.line), func(t *testing.T) {
				b1.Write(test.input)
				b2.WriteString(string(test.input))

				if diff := cmp.Diff(test.bytes, b1.Bytes()); diff != "" {
					t.Errorf("b1: incorrect result (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(test.bytes, b2.Bytes()); diff != "" {
					t.Errorf("b2: incorrect result (-want +got):\n%s", diff)
				}

				gotLen := b1.Len()
				if test.len != gotLen {
					t.Errorf("b1: incorrect len: want %d, got %d", test.len, gotLen)
				}
				gotLen = b2.Len()
				if test.len != gotLen {
					t.Errorf("b2: incorrect len: want %d, got %d", test.len, gotLen)
				}
			})
		}
	}
}

var out1 []byte

func FuzzRingBuf(f *testing.F) {
	var out []byte
	f.Fuzz(func(t *testing.T, size uint, write1 []byte, write2 []byte) {
		rb := ringbuf.New(int(size))
		for range 10 {
			rb.Write(write1)
			out = rb.Bytes()
			rb.Write(write2)
			out = rb.Bytes()
		}
	})
	out1 = out // avoid optimization
}

func BenchmarkRingBuf_Write(b *testing.B) {
	in := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	rb := ringbuf.New(10)
	b.ResetTimer()
	for b.Loop() {
		rb.Write(in)
	}
}

var out2 []byte

func BenchmarkRingBuf_Write_Bytes(b *testing.B) {
	in := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	out := []byte{}
	rb := ringbuf.New(10)
	b.ResetTimer()
	for b.Loop() {
		rb.Write(in)
		out = rb.Bytes()
	}
	out2 = out
}

func BenchmarkRingBuf_New_Write(b *testing.B) {
	bytes := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	b.ResetTimer()
	for b.Loop() {
		rb := ringbuf.New(10)
		rb.Write(bytes)
	}
}

// line returns the line number where it's called.
func line() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}
