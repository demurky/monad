// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"errors"
	"os"
	"syscall"
)

func IntKill(h Handle) error {
	err := Interrupt(h)
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		err = Kill(h)
	}
	return err
}

func TermKill(h Handle) error {
	err := Signal(h, syscall.SIGTERM)
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		err = Kill(h)
	}
	return err
}

func SigKill(h Handle, sig os.Signal) error {
	err := Signal(h, sig)
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		err = Kill(h)
	}
	return err
}

func DescIntKill(h Handle) (err error) {
	desc, err := Descendents(h)
	if err != nil {
		return IntKill(h)
	}
	for _, v := range desc {
		err = IntKill(v)
	}
	return err
}

func DescTermKill(h Handle) (err error) {
	desc, err := Descendents(h)
	if err != nil {
		return TermKill(h)
	}
	for _, v := range desc {
		err = TermKill(v)
	}
	return err
}

func DescSigKill(h Handle, sig os.Signal) (err error) {
	desc, err := Descendents(h)
	if err != nil {
		return SigKill(h, sig)
	}
	for _, v := range desc {
		err = SigKill(v, sig)
	}
	return err
}

func Interrupt(h Handle) error {
	return interrupt(h)
}

func Kill(h Handle) error {
	return kill(h)
}

func Signal(h Handle, sig os.Signal) error {
	return signal(h, sig)
}
