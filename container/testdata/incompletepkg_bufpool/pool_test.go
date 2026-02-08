// Copyright 2025 the toolbox authors.
// SPDX-License-Identifier: Apache-2.0

package bufpool_test

import (
	"math/rand/v2"
	"slices"
	"sync"
	"testing"

	"github.com/layer8co/toolbox/container/bufpool"
)

func BenchmarkGetPut(b *testing.B) {

	b.Run("bufpool", func(b *testing.B) {
		pool := bufpool.New[byte]()
		user := pool.NewUser()
		for b.Loop() {
			buf := user.Get()
			*buf = append(*buf, "ayylmao"...)
			user.Put(buf, len(*buf))
		}
	})

	b.Run("stdlib", func(b *testing.B) {
		pool := sync.Pool{
			New: func() any {
				b := make([]byte, 0, 16)
				return &b
			},
		}
		for b.Loop() {
			buf := pool.Get().(*[]byte)
			*buf = append(*buf, "ayylmao"...)
			pool.Put(buf)
		}
	})
}

func BenchmarkGetPutRising(b *testing.B) {

	b.Run("bufpool", func(b *testing.B) {
		pool := bufpool.New[byte]()
		user := pool.NewUser()
		i := uint16(0)
		for b.Loop() {
			buf := user.Get()
			*buf = slices.Grow(*buf, int(i))
			user.Put(buf, len(*buf))
			i++
		}
	})

	b.Run("stdlib", func(b *testing.B) {
		pool := sync.Pool{
			New: func() any {
				b := make([]byte, 0, 16)
				return &b
			},
		}
		i := uint16(0)
		for b.Loop() {
			buf := pool.Get().(*[]byte)
			*buf = slices.Grow(*buf, int(i))
			pool.Put(buf)
			i++
		}
	})
}

func BenchmarkGetPutRand(b *testing.B) {

	x := 100
	y := 100_000

	b.Run("bufpool", func(b *testing.B) {
		pool := bufpool.New[byte]()
		user := pool.NewUser()
		for b.Loop() {
			buf := user.Get()
			*buf = slices.Grow(*buf, randBetween(x, y))
			user.Put(buf, len(*buf))
		}
	})

	b.Run("stdlib", func(b *testing.B) {
		pool := sync.Pool{
			New: func() any {
				b := make([]byte, 0, 16)
				return &b
			},
		}
		for b.Loop() {
			buf := pool.Get().(*[]byte)
			*buf = slices.Grow(*buf, randBetween(x, y))
			pool.Put(buf)
		}
	})
}

// randBetween returns random number n such that a <= n <= b.
func randBetween(a, b int) int {
	if a > b {
		panic("a > b")
	}
	return rand.IntN(b-a+1) + a
}

// func BenchmarkPoolCap(b *testing.B) {
//
// 	pool := bufpool.New[byte]()
// 	expectedCap := 512
//
// 	for b.Loop() {
// 		buf := pool.Get(expectedCap)
// 		*buf.Buf = append(*buf.Buf, "ayylmao"...)
// 		buf.Put()
// 	}
// }
//
// func BenchmarkPoolId(b *testing.B) {
//
// 	pool := bufpool.New[byte]()
// 	uniq := new(byte)
//
// 	for b.Loop() {
// 		buf := pool.Get(uniq)
// 		*buf.Buf = append(*buf.Buf, "wut"...)
// 		buf.Put()
// 	}
// }
