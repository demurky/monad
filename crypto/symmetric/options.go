// Copyright 2025 the toolbox authors.
// SPDX-License-Identifier: Apache-2.0

package symmetric

import "fmt"

type Mode uint8

const (
	modeBegin Mode = iota
	////
	ModeXChaCha20
	ModeAES256CTR
	////
	modeEnd
)

func (m Mode) String() string {
	switch m {
	case ModeXChaCha20:
		return "XChaCha20"
	case ModeAES256CTR:
		return "AES256-CTR"
	default:
		panic(fmt.Sprintf("symmetric: unknown mode %d", m))
	}
}

type config struct {
	mode Mode

	argonTime    uint32
	argonTimeMax uint32

	argonMemory    uint32
	argonMemoryMax uint32

	argonThreads    uint8
	argonThreadsMax uint8
}

func getConfig(options []Option) *config {

	c := &config{

		// Make sure these values are reflected
		// in the documentations of [NewEncryptor] and [NewDecryptor].

		mode: ModeXChaCha20,

		argonTime:    3,
		argonTimeMax: 10,

		argonMemory:    16 * 1024,
		argonMemoryMax: 64 * 1024,

		argonThreads:    8,
		argonThreadsMax: 64,
	}

	for _, fn := range options {
		fn(c)
	}

	return c
}

type Option func(*config)

func WithMode(mode Mode) Option {
	return func(c *config) {
		c.mode = mode
	}
}

func WithArgonTime(time uint32) Option {
	return func(c *config) {
		c.argonTime = time
	}
}

func WithArgonTimeMax(time uint32) Option {
	return func(c *config) {
		c.argonTimeMax = time
	}
}

func WithArgonMemory(mem uint32) Option {
	return func(c *config) {
		c.argonMemory = mem
	}
}

func WithArgonMemoryMax(mem uint32) Option {
	return func(c *config) {
		c.argonMemoryMax = mem
	}
}

func WithArgonThreads(n uint8) Option {
	return func(c *config) {
		c.argonThreads = n
	}
}

func WithArgonThreadsMax(n uint8) Option {
	return func(c *config) {
		c.argonThreadsMax = n
	}
}
