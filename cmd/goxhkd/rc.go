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
