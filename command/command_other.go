// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package command

import (
	"os"
	"os/exec"
)

func interrupt(h Handle) error {
	return signal(h, os.Interrupt)
}

func kill(h Handle) error {
	return signal(h, os.Kill)
}

func signal(h Handle, sig os.Signal) error {
	p, err := toProcess(h)
	if err != nil {
		return err
	}
	return p.Signal(sig)
}

func windowsNewProcessGroup(cmd *exec.Cmd) {}
