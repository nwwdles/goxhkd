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
		xdgConfigHome = path.Join(home, ".local", "share")
	}

	err := tryExecRcPath(path.Join(xdgConfigHome, rcSubdir, rcName))
	if err == nil {
		return nil // return on success
	}

	err = tryExecRcPath(path.Join(home, "."+rcName))
	if err == nil {
		return nil // return on success
	}

	return err
}
