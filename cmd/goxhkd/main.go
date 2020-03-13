package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/cupnoodles14/scratchpad/go/goxhkd/pkg/comm"
)

type RPCAdapter struct {
	X *xgbutil.XUtil
}

func (r *RPCAdapter) BindCommand(binding comm.Binding, result *bool) error {
	bindCommand(r.X, binding.Btn, binding.Cmd, binding.RunOnPress, binding.Repeating)
	return nil
}

func (r *RPCAdapter) ClearAll(*bool, *bool) error {
	unbindAll(r.X)
	return nil
}

func ListenAndServe() error {
	c := comm.Connection{
		Network: "unix",
		Address: comm.SocketAddr,
	}

	if c.Network == "unix" {
		if err := os.RemoveAll(c.Address); err != nil {
			log.Fatal(err)
		}
	}

	ln, err := net.Listen(c.Network, c.Address)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(c)
	}
}

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	ra := &RPCAdapter{
		X: X,
	}

	err = rpc.Register(ra)
	if err != nil {
		log.Fatalf(err.Error())
	}

	keybind.Initialize(X)

	// bindCommand(X, "Mod4-v", `notify-send 'hello world!'`, false, false)
	// bindCommand(X, "Mod4-shift-v", `notify-send 'hello world 2!'`, false, true)

	serverErrors := make(chan error, 1)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() { serverErrors <- ListenAndServe() }()
	go func() { xevent.Main(X) }()

	log.Println("Program initialized. Start pressing keys!")
	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatal(err)
		}
	case sig := <-osSignals:
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			return
		}
	}
}
