// Copyright 2025 the toolbox authors.
// SPDX-License-Identifier: Apache-2.0

package streamcrypt

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBadChecksum(t *testing.T) {

	testingBadChecksum = true
	defer func() {
		testingBadChecksum = false
	}()

	tests := []struct {
		mode Mode
	}{
		{mode: ModeXChaCha20},
		{mode: ModeAES256CTR},
	}

	input := []byte("hello world")

	passwordString := "mypass123"
	password := []byte(passwordString)
	passFunc := func() ([]byte, error) {
		return []byte(passwordString), nil
	}

	for _, test := range tests {
		t.Run(test.mode.String(), func(t *testing.T) {

			ciphertext := Encrypt(input, password, WithMode(test.mode))
			output, err := Decrypt(ciphertext, passFunc)

			if !errors.Is(err, ErrBadChecksum) {
				t.Errorf("incorrect error: want ErrBadChecksum, got error %q", err)
			}

			if diff := cmp.Diff(input, output); diff != "" {
				t.Errorf("incorrect result (-want +got):\n%s", diff)
			}
		})
	}
}
