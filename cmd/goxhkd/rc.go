// The MIT License (MIT)
//
// Copyright (c) 2020 cupnoodles
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"errors"
	"os"
	"os/exec"
	"path"
)

const (
	rcEnv    = "GOXHKD_RC"
	rcName   = "goxhkdrc"
	rcSubdir = appName
)

func getSearchPaths() []string {
	filepath := os.Getenv(rcEnv)
	if filepath != "" {
		return []string{filepath}
	}

	home := os.Getenv("HOME")

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = path.Join(home, ".config")
	}

	return []string{
		path.Join(xdgConfigHome, rcSubdir, rcName),
		path.Join(home, "."+rcName),
		"/etc/goxhkd/goxhkdrc",
	}
}

func runRc() error {
	for _, path := range getSearchPaths() {
		err := exec.Command(path).Run()

		switch {
		case err == nil:
			return nil
		case !errors.Is(err, os.ErrNotExist):
			return err
		}
	}

	return nil
}
