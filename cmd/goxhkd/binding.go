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
)

func keyIsPressed(x *xgbutil.XUtil, keycode xproto.Keycode) (bool, error) {
	ck := xproto.QueryKeymap(x.Conn())

	reply, err := ck.Reply()
	if err != nil {
		err = fmt.Errorf("failed to get key press state: %w", err)

		return false, err
	}

	return reply.Keys[keycode>>3]&(0x1<<(keycode%8)) != 0, err
}

func bindCommand(x *xgbutil.XUtil, w xproto.Window, btn string, cmd []string, runOnRelease, repeating, sh bool) error {
	log.Printf("grabbing %s (window: %d): %s", btn, w, cmd)
	runner := func() error {
		if sh {
			cmd = append([]string{"sh", "-c"}, cmd...)
		}

		log.Println("Running:", cmd)

		return exec.Command(cmd[0], cmd[1:]...).Start() // #nosec
	}

	if repeating {
		return bindCmdRepeating(x, w, btn, runOnRelease, runner)
	}

	return bindCmd(x, w, btn, runOnRelease, runner)
}

func bindCmdRepeating(x *xgbutil.XUtil, w xproto.Window, btn string, runOnRelease bool, runCmd func() error) error {
	if !runOnRelease {
		return keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
			if err := runCmd(); err != nil {
				log.Println(err)
			}
		}).Connect(x, w, btn, true)
	}

	return keybind.KeyReleaseFun(func(x *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		if err := runCmd(); err != nil {
			log.Println(err)
		}
	}).Connect(x, w, btn, true)
}

func bindCmd(x *xgbutil.XUtil, w xproto.Window, btn string, runOnRelease bool, runner func() error) error {
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
			p, err := keyIsPressed(x, e.Keycode())
			if err != nil {
				log.Print(err)

				return
			}

			if runOnRelease && p {
				return
			}

			if err = runner(); err != nil {
				log.Print(err)
			}
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
	}).Connect(x, w, btn, true)
	if err != nil {
		return err
	}

	return keybind.KeyReleaseFun(func(x *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		releaseFun(KeyReleaseEvent(e))
	}).Connect(x, w, btn, true)
}
