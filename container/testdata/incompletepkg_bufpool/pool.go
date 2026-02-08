// Copyright 2025 the toolbox authors.
// SPDX-License-Identifier: Apache-2.0

// Package bufpool implements a dynamic slice pool.
package bufpool

import (
	"math/bits"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// In the worst case, if you set this value to e.g. 10 seconds,
	// the maximum rate of allocations will be limited to about
	// one allocation every 10 seconds.
	estimationPeriod = 10 * time.Second

	// The number of periods a user needs to be idle for
	// to be considered idle.
	userIdlePeriods = 6

	// Make the first bucket be for sizes less than or equal 2^minBucketBits.
	minBucketBits = 4

	minBucketSize = 1 << minBucketBits
	bucketsCount  = strconv.IntSize - minBucketBits
)

var testCleanup func() = nil

type Pool[T any] struct {
	buckets  [bucketsCount]sync.Pool
	ticker   *time.Ticker
	mu       *sync.Mutex
	users    map[*PoolUser[T]]struct{}
	userPool sync.Pool
}

type PoolUser[T any] struct {
	pool             *Pool[T]
	sizeEstimate     atomic.Int64
	sizeEstimateNext atomic.Int64
	registered       atomic.Bool
	idlePeriods      int
}

func New[T any]() *Pool[T] {

	ticker := time.NewTicker(estimationPeriod)
	ticker.Stop()

	mu := new(sync.Mutex)
	users := make(map[*PoolUser[T]]struct{})

	p := &Pool[T]{
		ticker: ticker,
		mu:     mu,
		users:  users,
		userPool: sync.Pool{
			New: func() any {
				return &PoolUser[T]{}
			},
		},
	}
	for i := range p.buckets {
		p.buckets[i].New = func() any {
			b := make([]T, 0, bucketSize(i))
			return &b
		}
	}

	tick := func() {

		mu.Lock()
		defer mu.Unlock()

		for u := range users {

			next := u.sizeEstimateNext.Swap(0)
			u.sizeEstimate.Store(next)

			if isIdle := next == 0; !isIdle {
				u.idlePeriods = 0
				continue
			}

			u.idlePeriods++

			if u.idlePeriods >= userIdlePeriods {
				u.idlePeriods = 0
				delete(users, u)
				u.registered.Store(false)
			}
		}

		if len(users) == 0 {
			ticker.Stop()
		}
	}

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				tick()
			case <-done:
				if testCleanup != nil {
					testCleanup()
				}
				return
			}
		}
	}()

	runtime.AddCleanup(p, func(ch chan struct{}) { close(ch) }, done)

	return p
}

// Calling [PoolUser.Done]
// merely saves on allocations and is not necessary.
func (p *Pool[T]) NewUser() *PoolUser[T] {
	u := p.userPool.Get().(*PoolUser[T])
	u.pool = p
	return u
}

func (p *Pool[T]) get(size int) *[]T {
	return p.buckets[bucketIndex(size)].Get().(*[]T)
}

func (p *Pool[T]) put(b *[]T) {
	p.buckets[bucketIndex(cap(*b))].Put(b)
}

// ==================================================

func (u *PoolUser[T]) Get(size ...int) *[]T {
	if len(size) > 0 {
		if size[0] < 0 {
			panic("bufpool.PoolUser.Get: size[0] < 0")
		}
		return u.pool.get(size[0])
	}
	return u.pool.get(int(u.sizeEstimate.Load()))
}

// A negative size bypasses updating the size estimation values.
func (u *PoolUser[T]) Put(b *[]T, size int) {
	if size >= 0 {
		updateIfIncrease(&u.sizeEstimate, int64(size))
		updateIfIncrease(&u.sizeEstimateNext, int64(size))
		if !u.registered.Load() {
			u.register()
		}
	}
	*b = (*b)[:0]
	u.pool.put(b)
}

func (u *PoolUser[T]) Reset() {
	// Because of the lock-free way we update these atomics
	// in the background goroutine, it's important here
	// that we zero sizeEstimateNext before sizeEstimate.
	u.sizeEstimateNext.Store(0)
	u.sizeEstimate.Store(0)
}

func (u *PoolUser[T]) Done() {

	p := u.pool

	if p == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.users, u)
	if len(p.users) == 0 {
		p.ticker.Stop()
	}

	*u = PoolUser[T]{}
	p.userPool.Put(u)
}

func (u *PoolUser[T]) register() {

	p := u.pool

	p.mu.Lock()
	defer p.mu.Unlock()

	if u.registered.Load() {
		return
	}

	p.users[u] = struct{}{}
	u.registered.Store(true)
	if len(p.users) == 1 {
		p.ticker.Reset(estimationPeriod)
	}
}

// ==================================================

func bucketIndex(size int) int {
	if size <= minBucketSize {
		return 0
	}
	return bits.Len(uint(size-1)) - minBucketBits
}

func bucketSize(index int) int {
	// TODO: This overflows on (index + minBucketBits) == 0;
	// it should be `(1 << (index + minBucketSize)) - 1`,
	// but then bucketIndex would need adjustments to
	// to correctly report the index of each given size.
	// At the end of the day, we should make bucketIndex(math.MaxInt)
	// possible without overflows.
	return 1 << (index + minBucketBits)
}

// updateIfIncrease atomically updates the given atomic.Int64
// if new is larger than it.
func updateIfIncrease(a *atomic.Int64, new int64) {
	// Atomically compare and update the atomic using a CAS loop.
	for {
		old := a.Load()
		if new <= old {
			return
		}
		if a.CompareAndSwap(old, new) {
			return
		}
	}
}
