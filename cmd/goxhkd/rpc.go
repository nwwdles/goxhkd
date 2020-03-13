package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"

	"github.com/BurntSushi/xgbutil"
	"github.com/cupnoodles14/scratchpad/go/goxhkd/pkg/comm"
)

type GoRPC struct {
	X    *xgbutil.XUtil
	Conn *comm.Connection
}

func (r *GoRPC) listenAndServe() error {
	err := rpc.Register(r)
	if err != nil {
		return err
	}

	if r.Conn.Network == "unix" {
		if err = os.RemoveAll(r.Conn.Address); err != nil {
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

func (r *GoRPC) BindCommand(binding comm.Binding, _ *bool) error {
	return bindCommand(r.X, binding.Btn, binding.Cmd, binding.RunOnPress, binding.Repeating)
}

func (r *GoRPC) UnbindAll(_, _ *bool) error {
	return unbindAll(r.X)
}
