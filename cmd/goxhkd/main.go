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

const AppName = "goxhkd"
const InitialRcRunDelay = 200 * time.Millisecond

func main() {
	conn := shared.DefaultSocketConnection()
	flag.StringVar(&conn.Network, "network", conn.Network, "specify connection network (unix, tcp, ...)")
	flag.StringVar(&conn.Address, "address", conn.Address, "specify connection address (socket path, host, ...)")
	flag.Parse()

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(X)

	ra := GoRPC{
		X:    X,
		Conn: conn,
	}

	serverErrors := make(chan error, 1)

	go func() { serverErrors <- ra.listenAndServe() }()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() { xevent.Main(X) }()

	log.Println(AppName, "started. Start pressing keys!")

	// run RC file
	go func() {
		time.Sleep(InitialRcRunDelay)

		if err := runRc(); err != nil {
			fmt.Println("RC file couldn't be executed:", err)
		}
	}()

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
