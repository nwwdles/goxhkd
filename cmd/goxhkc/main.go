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
	"net/rpc"
	"os"

	"gitlab.com/cupnoodles14/goxhkd/pkg/shared"
)

// OS return codes
const (
	Success      = 0
	GenericError = 1
)

var (
	ErrNoAction = errors.New("button requires either -command or -clear action")
)

func run() int {
	conn := shared.DefaultSocketConnection()

	btn := flag.String("button", "", "specify a button")
	cmd := flag.String("command", "", "set command for the button")
	sh := flag.Bool("sh", false, "run command with 'sh -c ...'")
	onRelease := flag.Bool("onrelease", false, "run command on button release")
	repeating := flag.Bool("repeat", false, "repeatedly run command while the button is pressed")
	clear := flag.Bool("clear", false, "clear the button")
	clearAll := flag.Bool("clearall", false, "clear all bindings")

	flag.StringVar(&conn.Network, "network", conn.Network, "specify connection network (unix, tcp, ...)")
	flag.StringVar(&conn.Address, "address", conn.Address, "specify connection address (socket path, host, ...)")

	flag.Parse()

	c, err := rpc.Dial(conn.Network, conn.Address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return GenericError
	}
	defer c.Close()

	if *btn != "" {
		switch {
		case *clear:
			err = c.Call("GoRPC.Unbind", shared.Binding{
				Btn:          *btn,
				RunOnRelease: *onRelease,
			}, nil)
		case *cmd != "":
			err = c.Call("GoRPC.BindCommand", shared.Binding{
				Cmd:          *cmd,
				Btn:          *btn,
				RunOnRelease: *onRelease,
				Repeating:    *repeating,
				Sh:           *sh,
			}, nil)
		default:
			err = ErrNoAction
		}
	} else if *clearAll {
		err = c.Call("GoRPC.UnbindAll", struct{}{}, nil)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return GenericError
	}

	return Success
}

func main() {
	os.Exit(run())
}
