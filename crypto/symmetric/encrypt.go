// Copyright 2025 the github.com/demurky/monad authors.
// SPDX-License-Identifier: Apache-2.0

package symmetric

import (
	"crypto/cipher"
	"crypto/sha3"
	"io"
	"io/fs"

	"golang.org/x/crypto/argon2"
)

// Encryptor is returned by [NewEncryptor].
// See it's documentation for details.
//
// Encryptor implements [io.WriteCloser].
type Encryptor struct {
	dest          io.Writer
	stream        cipher.Stream
	header        header
	hash          *sha3.SHAKE
	firstTime     bool
	done          bool
	ciphertextBuf []byte // Buffer used for encryption.
}

// NewEncryptor returns an [Encryptor]
// which is an [io.WriteCloser]
// that encrypts plaintext and writes the ciphertext to dest.
//
// [Encryptor.Close] must be called after all writes are concluded
// in order to write the authentication bytes to dest.
//
// The password is not retained by this function.
//
// The following options can be used to configure the encryption behavior:
//   - [WithMode] (default: [ModeXChaCha20])
//   - [WithArgonTime] (default: 3)
//   - [WithArgonMemory] (default: 16*1024)
//   - [WithArgonThreads] (default: 8)
func NewEncryptor(
	dest io.Writer,
	password []byte,
	options ...Option,
) *Encryptor {

	e := &Encryptor{
		dest:      dest,
		firstTime: true,
		header:    newHeader(getConfig(options)),
	}

	key := argon2.IDKey(
		password,
		e.header.ArgonSalt[:],
		e.header.ArgonTime,
		e.header.ArgonMemory,
		e.header.ArgonThreads,
		e.header.keyLen(),
	)

	e.stream = e.header.getStream(key)
	e.hash = sha3.NewSHAKE256()
	e.hash.Write(key)
	clear(key)

	e.header.writeTo(e.hash)

	return e
}

func (e *Encryptor) Write(plaintext []byte) (int, error) {

	if e.done {
		return 0, fs.ErrClosed
	}

	err := e.writeHeader()
	if err != nil {
		return 0, err
	}

	if cap(e.ciphertextBuf) < len(plaintext) {
		e.ciphertextBuf = make([]byte, len(plaintext))
	}
	e.ciphertextBuf = e.ciphertextBuf[:len(plaintext)]

	e.stream.XORKeyStream(e.ciphertextBuf, plaintext)
	e.hash.Write(e.ciphertextBuf)
	return e.dest.Write(e.ciphertextBuf)
}

func (e *Encryptor) Close() error {
	e.done = true
	err := e.writeHeader()
	if err != nil {
		return err
	}
	_, err = e.dest.Write(getChecksum(e.hash))
	return err
}

func (e *Encryptor) writeHeader() error {
	if e.firstTime {
		e.firstTime = false
		err := e.header.writeTo(e.dest)
		if err != nil {
			return err
		}
	}
	return nil
}
