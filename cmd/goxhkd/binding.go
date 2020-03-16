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
	"fmt"
	"log"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/google/shlex"
)

func makeCommand(cmd string, sh bool) (*exec.Cmd, error) {
	if sh {
		return exec.Command("sh", "-c", cmd), nil
	}

	tokens, err := shlex.Split(cmd)
	if err != nil {
		return nil, err
	}

	return exec.Command(tokens[0], tokens[1:]...), nil
}

func makeCmdRunner(cmd string, sh bool) func() error {
	return func() error {
		log.Println("Running:", cmd)

		cmd, err := makeCommand(cmd, sh)
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

func bindCommand(x *xgbutil.XUtil, btn, cmd string, runOnRelease, repeating, sh bool) error {
	runner := makeCmdRunner(cmd, sh)

	if repeating {
		return bindCommandRepeating(x, btn, runOnRelease, runner)
	}

	return bindCommandNonrepeating(x, btn, runOnRelease, runner)
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func bindCommandRepeating(x *xgbutil.XUtil, btn string, runOnRelease bool, runner func() error) error {
	var err error

	if !runOnRelease {
		err = keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
			logErr(runner())
		}).Connect(x, x.RootWin(), btn, true)
	} else {
		err = keybind.KeyReleaseFun(func(x *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			logErr(runner())
		}).Connect(x, x.RootWin(), btn, true)
	}

	return err
}

func bindCommandNonrepeating(x *xgbutil.XUtil, btn string, runOnRelease bool, runner func() error) error {
	var (
		pressFun,
		releaseFun func(e timedKeyEvent)
		lastEventTime xproto.Timestamp // this variable is captured
	)

	// lastEventTime is used to filter out artificial events spawned by key
	// autorepeating. Such events come in pressEvent-releaseEvent pairs that
	// have the same timestamp
	timer := func(e timedKeyEvent) { lastEventTime = e.Timestamp() }

	executor := func(e timedKeyEvent) {
		t := e.Timestamp()
		if t != lastEventTime {
			lastEventTime = t

			// keyIsPressed is used to detect artificial events in cases when
			// the command is bound to key release.
			if runOnRelease && keyIsPressed(x, e.Keycode()) {
				return
			}

			logErr(runner())
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
