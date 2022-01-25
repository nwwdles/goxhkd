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
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/BurntSushi/xgbutil"
	"gitlab.com/cupnoodles14/goxhkd/pkg/shared"
)

type App struct {
	X    *xgbutil.XUtil
	Conn *shared.Connection
}

func (r *App) listenAndServe() (err error) {
	err = rpc.Register(r)
	if err != nil {
		return err
	}

	if r.Conn.Network == "unix" {
		err = os.RemoveAll(r.Conn.Address)
		if err != nil {
			return err
		}
	}

	ln, err := net.Listen(r.Conn.Network, r.Conn.Address)
	if err != nil {
		return err
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go rpc.ServeConn(c)
	}
}

func (r *App) BindCommand(b shared.Binding, _ *struct{}) error {
	w := b.Window
	if w == 0 {
		w = r.X.RootWin()
	}

	return bindCommand(r.X, w, b.Btn, b.Cmd, b.RunOnRelease, b.Repeating, b.Sh)
}

func (r *App) UnbindAll(_ struct{}, _ *struct{}) error {
	return unbindAll(r.X)
}

func (r *App) Unbind(b shared.Binding, _ *struct{}) error {
	// Because workaround for xorg key repeating uses both release and press
	// functions, we can't easily unbind just the release or just the press
	// event, so we unbind both.

	err := unbind(r.X, b.Btn, !b.RunOnRelease)
	if err != nil {
		log.Println(err)
	}

	return unbind(r.X, b.Btn, b.RunOnRelease)
}
