package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"

	"github.com/BurntSushi/xgbutil"
	"github.com/cupnoodles14/scratchpad/go/goxhkd/pkg/shared"
)

// GoRPC implements RPC adapter
type GoRPC struct {
	X    *xgbutil.XUtil
	Conn *shared.Connection
}

func (r *GoRPC) listenAndServe() error {
	err := rpc.Register(r)
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
		fmt.Println(err)
		return err
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go rpc.ServeConn(c)
	}
}

func (r *GoRPC) BindCommand(b shared.Binding, _ *struct{}) error {
	return bindCommand(r.X, b.Btn, b.Cmd, b.RunOnRelease, b.Repeating, b.Sh)
}

func (r *GoRPC) UnbindAll(_ struct{}, _ *struct{}) error {
	return unbindAll(r.X)
}

func (r *GoRPC) Unbind(b shared.Binding, _ *struct{}) error {
	// Because workaround for xorg key repeating uses both release and press
	// functions, we can't easily unbind just the release or just the press
	// event, so we unbind both.
	err := unbind(r.X, b.Btn, !b.RunOnRelease)
	if err != nil {
		_ = err // skip
	}

	return unbind(r.X, b.Btn, b.RunOnRelease)
}
