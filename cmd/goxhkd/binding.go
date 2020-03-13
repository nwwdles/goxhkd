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

func keyIsPressed(X *xgbutil.XUtil, keycode xproto.Keycode) bool {
	ck := xproto.QueryKeymap(X.Conn())
	reply, err := ck.Reply()

	if err != nil {
		err = fmt.Errorf(": %w", err)
		log.Printf("%v", err)
	}

	return reply.Keys[keycode>>3]&(0x1<<(keycode%8)) != 0
}

func bindCommand(X *xgbutil.XUtil, btn, cmd string, runOnPress, repeating bool) {
	if repeating {
		bindCommandRepeating(X, btn, cmd, runOnPress)
	} else {
		bindCommandNonrepeating(X, btn, cmd, runOnPress)
	}
}

func bindCommandRepeating(X *xgbutil.XUtil, btn, cmd string, runOnPress bool) error {
	runCmd := makeCmdRunner(cmd)

	var keyFun xevent.KeyPressFun
	if runOnPress {
		keyFun = keybind.KeyPressFun(
			func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
				runCmd()
			})
	} else {
		keyFun = keybind.KeyReleaseFun(
			func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
				runCmd()
			})
	}

	err := keyFun.Connect(X, X.RootWin(), btn, true)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func bindCommandNonrepeating(X *xgbutil.XUtil, btn, cmd string, runOnPress bool) {
	runCmd := makeCmdRunner(cmd)

	var (
		pressFun, releaseFun func(e timedKeyEvent)
		lastEventTime        xproto.Timestamp // this variable is captured
	)

	// lastEventTime is used to filter out artificial events spawned by key
	// autorepeating. Such events come in pressEvent-releaseEvent pairs that
	// have the same timestamp
	timer := func(e timedKeyEvent) { lastEventTime = e.GetTime() }

	executor := func(e timedKeyEvent) {
		t := e.GetTime()
		if t != lastEventTime {
			lastEventTime = t

			if !runOnPress && keyIsPressed(X, e.GetKeycode()) {
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

	err := keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
		pressFun(KeyPressEvent(e))
	}).Connect(X, X.RootWin(), btn, true)
	if err != nil {
		log.Fatal(err)
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		releaseFun(KeyReleaseEvent(e))
	}).Connect(X, X.RootWin(), btn, true)
	if err != nil {
		log.Fatal(err)
	}
}

func unbindAll(X *xgbutil.XUtil) {
	keybind.Detach(X, X.RootWin())
}
