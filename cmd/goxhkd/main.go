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
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"gitlab.com/cupnoodles14/goxhkd/pkg/shared"
)

const (
	appName           = "goxhkd"
	initialRcRunDelay = 200 * time.Millisecond
)

// set via ldflags
var (
	version = "dev"
	build   = ""
)

func main() {
	conn := shared.DefaultSocketConnection()
	flag.StringVar(&conn.Network, "network", conn.Network, "specify connection network (unix, tcp, ...)")
	flag.StringVar(&conn.Address, "address", conn.Address, "specify connection address (socket path, host, ...)")
	v := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *v {
		fmt.Println(version, build)
		return
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(X)

	ra := App{
		X:    X,
		Conn: conn,
	}

	serverErrors := make(chan error, 1)
	go func() { serverErrors <- ra.listenAndServe() }()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() { xevent.Main(X) }()

	// run RC file
	go func() {
		time.Sleep(initialRcRunDelay)

		if err := runRc(); err != nil {
			fmt.Println("RC file couldn't be executed:", err)
		}
	}()

	log.Println(appName, "started. Start pressing keys!")

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
