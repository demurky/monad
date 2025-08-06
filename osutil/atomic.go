// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package osutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {

	f, err := OpenFileAtomic(path, os.O_WRONLY, perm)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		os.Remove(f.Name())
		return fmt.Errorf("could not write to file %q: %w", f.Name(), err)
	}

	return f.Close()
}

func OpenFileAtomic(path string, flag int, perm os.FileMode) (File, error) {

	f, err := CreateTemp(path, os.O_WRONLY, perm)
	if err != nil {
		return File{}, fmt.Errorf(
			"could not create temporary file in %q: %w",
			filepath.Dir(path), err,
		)
	}

	close := func() error {
		err := f.Close()
		if err != nil {
			return fmt.Errorf("could not close file %q: %w", f.Name(), err)
		}
		err = os.Rename(f.Name(), path)
		if err != nil {
			return fmt.Errorf(
				"could not move file %q to %q: %w",
				f.Name(), path, err,
			)
		}
		return nil
	}

	return File{
		File:  f,
		close: close,
	}, nil
}

type File struct {
	*os.File
	close func() error
}

func (f File) Close() error {
	return f.close()
}

func (f File) CloseFile() error {
	return f.File.Close()
}
