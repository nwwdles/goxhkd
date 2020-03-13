package main

import (
	"flag"
	"net/rpc"

	"github.com/cupnoodles14/scratchpad/go/goxhkd/pkg/shared"
)

func main() {
	conn := shared.DefaultSocketConnection()

	btn := flag.String("button", "", "specify a button")
	cmd := flag.String("command", "", "set command for the button")
	onRelease := flag.Bool("onrelease", false, "run command on button release")
	repeating := flag.Bool("repeat", false, "repeatedly run command while the button is pressed")
	clear := flag.Bool("clear", false, "clear the button")
	clearAll := flag.Bool("clearall", false, "clear all bindings")

	flag.StringVar(&conn.Network, "network", conn.Network, "specify connection network (unix, tcp, ...)")
	flag.StringVar(&conn.Address, "address", conn.Address, "specify connection address (socket path, host, ...)")

	flag.Parse()

	c, err := rpc.Dial(conn.Network, conn.Address)
	if err != nil {
		panic(err)
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
			}, nil)
		default:
			panic("button requires a -command or -clear action")
		}
	} else if *clearAll {
		err = c.Call("GoRPC.UnbindAll", struct{}{}, nil)
	}

	if err != nil {
		panic(err.Error())
	}
}
