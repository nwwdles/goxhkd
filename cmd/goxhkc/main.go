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
