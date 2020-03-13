package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/cupnoodles14/scratchpad/go/goxhkd/pkg/shared"
)

const AppName = "goxhkd"

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(X)

	ra := GoRPC{
		X: X,
		Conn: &shared.Connection{
			Network: "unix",
			Address: shared.DefaultSocketAddr,
		},
	}

	serverErrors := make(chan error, 1)

	go func() { serverErrors <- ra.listenAndServe() }()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() { xevent.Main(X) }()

	log.Println(AppName, "started. Start pressing keys!")

	// run RC file
	go func() {
		time.Sleep(500 * time.Millisecond)

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
