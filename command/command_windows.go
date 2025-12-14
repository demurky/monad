// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package command

import (
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func interrupt(h Handle) error {
	return windows.GenerateConsoleCtrlEvent(
		windows.CTRL_BREAK_EVENT,
		uint32(toPid(h)),
	)
}

func kill(h Handle) error {
	p, err := toProcess(h)
	if err != nil {
		return err
	}
	return p.Kill()
}

func signal(h Handle, sig os.Signal) error {
	p, err := toProcess(h)
	if err != nil {
		return err
	}
	return p.Signal(sig)
}

func windowsNewProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags |= windows.CREATE_NEW_PROCESS_GROUP
}
