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

func makeCmdRunner(cmd string) func() {
	return func() {
		log.Println("Key press! Running:", cmd)

		cmd, err := makeCommand(cmd)
		if err != nil {
			log.Println(err)
		}

		err = cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}
}

func keyIsPressed(x *xgbutil.XUtil, keycode xproto.Keycode) bool {
	ck := xproto.QueryKeymap(x.Conn())
	reply, err := ck.Reply()

	if err != nil {
		err = fmt.Errorf(": %w", err)
		log.Printf("%v", err)
	}

	return reply.Keys[keycode>>3]&(0x1<<(keycode%8)) != 0
}

func bindCommand(x *xgbutil.XUtil, btn, cmd string, runOnPress, repeating bool) error {
	if repeating {
		return bindCommandRepeating(x, btn, cmd, runOnPress)
	}

	return bindCommandNonrepeating(x, btn, cmd, runOnPress)
}

func bindCommandRepeating(x *xgbutil.XUtil, btn, cmd string, runOnPress bool) error {
	var err error

	runCmd := makeCmdRunner(cmd)

	if runOnPress {
		err = keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
			runCmd()
		}).Connect(x, x.RootWin(), btn, true)
	} else {
		err = keybind.KeyReleaseFun(func(x *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			runCmd()
		}).Connect(x, x.RootWin(), btn, true)
	}

	return err
}

func bindCommandNonrepeating(x *xgbutil.XUtil, btn, cmd string, runOnPress bool) error {
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
			if !runOnPress && keyIsPressed(x, e.GetKeycode()) {
				return
			}

			runCmd()
		}
	}

	if runOnPress {
		pressFun = executor
		releaseFun = timer
	} else {
		pressFun = timer
		releaseFun = executor
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

func unbindAll(x *xgbutil.XUtil) error {
	keybind.Detach(x, x.RootWin())
	return nil
}
