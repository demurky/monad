// Copyright 2025 the github.com/koonix/x authors.
// SPDX-License-Identifier: Apache-2.0

package osutil

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJson(
	path string,
	data any,
	perm os.FileMode,
	prefix, indent string,
) error {

	f, err := OpenFileAtomic(path, os.O_WRONLY, perm)
	if err != nil {
		return err
	}

	e := json.NewEncoder(f)
	e.SetIndent(prefix, indent)

	err = e.Encode(data)
	if err != nil {
		os.Remove(f.Name())
		return fmt.Errorf("could not write json to %q: %w", f.Name(), err)
	}

	return f.Close()
}

func ReadJson(path string, dest any) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(dest)
	if err != nil {
		return fmt.Errorf("could not read json from %q: %w", path, err)
	}

	return nil
}
