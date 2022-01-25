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
	"flag"
	"fmt"
	"log"
	"net/rpc"

	"github.com/BurntSushi/xgb/xproto"
	"gitlab.com/cupnoodles14/goxhkd/pkg/shared"
)

var ErrNoAction = errors.New("button requires either -command or -clear action")

// set via ldflags
var (
	version = "0.1.0"
	build   = ""
)

func run() (err error) {
	conn := shared.DefaultSocketConnection()
	flag.StringVar(&conn.Network, "network", conn.Network, "specify connection network (unix, tcp, ...)")
	flag.StringVar(&conn.Address, "address", conn.Address, "specify connection address (socket path, host, ...)")

	btn := flag.String("button", "", "specify a button")
	window := flag.Uint("window", 0, "specify a window")
	sh := flag.Bool("sh", false, "run command with 'sh -c ...'")
	onRelease := flag.Bool("onrelease", false, "run command on button release")
	repeating := flag.Bool("repeat", false, "repeatedly run command while the button is pressed")
	multi := flag.Bool("multi", false, "allow for multiple bindings to the same button")
	clearAll := flag.Bool("clearall", false, "clear all bindings")
	v := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *v {
		fmt.Println(version, build)
		return
	}

	c, err := rpc.Dial(conn.Network, conn.Address)
	if err != nil {
		return err
	}
	defer c.Close()

	switch {
	case *btn == "":
		if *clearAll {
			return c.Call("App.UnbindAll", struct{}{}, nil)
		}

		return ErrNoAction
	case len(flag.Args()) > 0:
		if !*multi {
			_ = c.Call("App.Unbind", shared.Binding{
				Btn:          *btn,
				RunOnRelease: *onRelease,
			}, nil)
		}
		return c.Call("App.BindCommand", shared.Binding{
			Cmd:          flag.Args(),
			Btn:          *btn,
			RunOnRelease: *onRelease,
			Repeating:    *repeating,
			Sh:           *sh,
			Window:       xproto.Window(*window),
		}, nil)
	default:
		return c.Call("App.Unbind", shared.Binding{
			Btn:          *btn,
			RunOnRelease: *onRelease,
		}, nil)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
