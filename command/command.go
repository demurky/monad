// Copyright 2025 the monad authors.
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"time"

	psutil "github.com/shirou/gopsutil/v4/process"
)

var ErrNotImplemented = errors.New("not implemented")

// Process represents an OS process and it's
type Process struct {
	*psutil.Process
	Level int
}

// Handle represents an OS process and may be any of the following types:
//
//   - int (indicatig the PID)
//   - int32 (indicating the PID)
//   - [*exec.Cmd]
//   - [*os.Process]
//   - [Process]
//   - [*psutil.Process]
//
//x:x
type Handle = any

// Graceful configures cmd so that
// it's signaled to shut down gracefully
// when its context is canceled.
// It's then forcefully killed if it's not dead after killTimeout.
func Graceful(
	cmd *exec.Cmd,
	killTimeout time.Duration,
	action ...func(h Handle) error,
) {
	WindowsNewProcessGroup(cmd)
	cmd.WaitDelay = killTimeout
	fn := IntKill
	if len(action) > 0 {
		fn = action[0]
	}
	cmd.Cancel = func() error {
		return fn(cmd)
	}
}

// Run runs cmd while capturing and it's outputs.
func Run(cmd *exec.Cmd) (stdout, stderr []byte, err error) {
	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf
	err = cmd.Run()
	if err != nil {
		return nil, nil, fmt.Errorf(
			"command %q failed: %w\nstdout:\n%s\nstderr:\n%s",
			cmd, err, stdoutBuf.Bytes(), stderrBuf.Bytes(),
		)
	}
	return stdoutBuf.Bytes(), stderrBuf.Bytes(), nil
}

// Shell returns the arguments for running a command in the system shell.
func Shell() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"cmd.exe", "/C"}
	case "plan9":
		return []string{"/bin/rc", "-c"}
	default:
		return []string{"/bin/sh", "-c"}
	}
}

// WindowsNewProcessGroup makes cmd's process be in a new process group
// when it's spawned on Windows. It's a no-op on other operating systems.
func WindowsNewProcessGroup(cmd *exec.Cmd) {
	windowsNewProcessGroup(cmd)
}

// Descendents returns the descendents of the given process,
// which excludes the process itself.
func Descendents(h Handle) ([]Process, error) {

	p, err := psutil.NewProcess(int32(toPid(h)))
	if err != nil {
		return nil, err
	}

	procs := make([]Process, 0)

	err = descendents(p, &procs, 0)
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(procs, func(a, b Process) int {
		return -cmp.Compare(a.Level, b.Level)
	})

	return procs, nil
}

func descendents(p *psutil.Process, procs *[]Process, level int) error {

	children, err := p.Children()
	if err != nil {
		return err
	}

	for _, child := range children {
		err := descendents(child, procs, level+1)
		if err != nil {
			continue
		}
	}

	*procs = append(*procs, Process{
		Process: p,
		Level:   level,
	})

	return nil
}

func toProcess(h Handle) (*os.Process, error) {
	switch h := h.(type) {
	case int:
		return os.FindProcess(h)
	case int32:
		return os.FindProcess(int(h))
	case *exec.Cmd:
		return h.Process, nil
	case *os.Process:
		return h, nil
	case Process:
		return os.FindProcess(int(h.Pid))
	case *psutil.Process:
		return os.FindProcess(int(h.Pid))
	default:
		panic(fmt.Sprintf("cmd.toProcess: unknown type %t", h))
	}
}

func toPid(h Handle) int {
	switch h := h.(type) {
	case int:
		return h
	case int32:
		return int(h)
	case *exec.Cmd:
		return h.Process.Pid
	case *os.Process:
		return h.Pid
	case Process:
		return int(h.Pid)
	case *psutil.Process:
		return int(h.Pid)
	default:
		panic(fmt.Sprintf("cmd.toPid: unknown type %t", h))
	}
}
