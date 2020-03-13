package main

import (
	"net/rpc"

	"github.com/cupnoodles14/scratchpad/go/goxhkd/pkg/comm"
)

func main() {
	conn := comm.Connection{
		Network: "unix",
		Address: comm.SocketAddr,
	}

	c, err := rpc.Dial(conn.Network, conn.Address)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	err = c.Call("RPCAdapter.BindCommand", comm.Binding{
		Cmd: "notify-send hello",
		Btn: "Mod4-v",
	}, nil)

	if err != nil {
		panic(err)
	}
}
