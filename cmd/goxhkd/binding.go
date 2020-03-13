package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/google/shlex"
)

func makeCommand(cmd string) (*exec.Cmd, error) {
	tokens, err := shlex.Split(cmd)
	return exec.Command(tokens[0], tokens[1:]...), err
}

func makeCmdRunner(cmd string) func() error {
	return func() error {
		log.Println("Key press! Running:", cmd)

		cmd, err := makeCommand(cmd)
		if err != nil {
			return err
		}

		return cmd.Start()
	}
}

func keyIsPressed(x *xgbutil.XUtil, keycode xproto.Keycode) bool {
	ck := xproto.QueryKeymap(x.Conn())
	reply, err := ck.Reply()

	if err != nil {
		err = fmt.Errorf("failed to get key press state: %w", err)
		log.Printf("%v", err)

		return false
	}

	return reply.Keys[keycode>>3]&(0x1<<(keycode%8)) != 0
}

func bindCommand(x *xgbutil.XUtil, btn, cmd string, runOnRelease, repeating bool) error {
	if repeating {
		return bindCommandRepeating(x, btn, cmd, runOnRelease) // TODO runOnRelease
	}

	return bindCommandNonrepeating(x, btn, cmd, runOnRelease)
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func bindCommandRepeating(x *xgbutil.XUtil, btn, cmd string, runOnRelease bool) error {
	var err error

	runCmd := makeCmdRunner(cmd)

	if !runOnRelease {
		err = keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
			logErr(runCmd())
		}).Connect(x, x.RootWin(), btn, true)
	} else {
		err = keybind.KeyReleaseFun(func(x *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			logErr(runCmd())
		}).Connect(x, x.RootWin(), btn, true)
	}

	return err
}

func bindCommandNonrepeating(x *xgbutil.XUtil, btn, cmd string, runOnRelease bool) error {
	runCmd := makeCmdRunner(cmd)

	var (
		pressFun,
		releaseFun func(e timedKeyEvent)
		lastEventTime xproto.Timestamp // this variable is captured
	)

	// lastEventTime is used to filter out artificial events spawned by key
	// autorepeating. Such events come in pressEvent-releaseEvent pairs that
	// have the same timestamp
	timer := func(e timedKeyEvent) { lastEventTime = e.GetTime() }

	executor := func(e timedKeyEvent) {
		t := e.GetTime()
		if t != lastEventTime {
			lastEventTime = t

			// keyIsPressed is used to detect artificial events in cases when
			// the command is bound to key release.
			if runOnRelease && keyIsPressed(x, e.GetKeycode()) {
				return
			}

			logErr(runCmd())
		}
	}

	if runOnRelease {
		pressFun = timer
		releaseFun = executor
	} else {
		pressFun = executor
		releaseFun = timer
	}

	err := keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
		pressFun(KeyPressEvent(e))
	}).Connect(x, x.RootWin(), btn, true)
	if err != nil {
		return err
	}

	err = keybind.KeyReleaseFun(func(x *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		releaseFun(KeyReleaseEvent(e))
	}).Connect(x, x.RootWin(), btn, true)

	return err
}
