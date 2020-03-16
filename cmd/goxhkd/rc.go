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
	"log"
	"os"
	"os/exec"
	"path"
)

const rcEnv = "GOXHKD_RC"
const rcName = "goxhkdrc"
const rcSubdir = AppName

func tryExecRcPath(filepath string) error {
	cmd := exec.Command(filepath)

	err := cmd.Run()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Println(err)
		}

		return err
	}

	return nil
}

func runRc() error {
	filepath := os.Getenv(rcEnv)
	if filepath != "" {
		return tryExecRcPath(filepath)
	}

	home := os.Getenv("HOME")

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = path.Join(home, ".config")
	}

	var err error
	for _, rc := range []string{
		path.Join(xdgConfigHome, rcSubdir, rcName),
		path.Join(home, "."+rcName),
		"/etc/goxhkd/goxhkdrc",
	} {
		err = tryExecRcPath(rc)
		if err == nil {
			return nil // return on success
		}
	}

	return err
}
