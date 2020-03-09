package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/google/shlex"
)

const defaultFIFOPath = "/tmp/goxhkd.fifo"

func makeCommand(cmd string) (*exec.Cmd, error) {
	tokens, err := shlex.Split(cmd)
	return exec.Command(tokens[0], tokens[1:]...), err
}

func runDaemon() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(X)

	bind := func(btn string, cmd string, onpress bool) {
		f := func() {
			log.Println("Key press! Running:", cmd)

			cmd, err := makeCommand(cmd)
			if err != nil {
				log.Println(err)
			}

			err = cmd.Start()
			if err != nil {
				log.Println(err)
			}
		}

		if onpress {
			var lastEventTime xproto.Timestamp

			kf := func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
				if e.Time != lastEventTime {
					lastEventTime = 0
					f()
					fmt.Println(e.Time)
				}
			}
			kr := func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
				lastEventTime = e.Time
			}

			err = keybind.KeyPressFun(kf).Connect(X, X.RootWin(), btn, true)
			err = keybind.KeyReleaseFun(kr).Connect(X, X.RootWin(), btn, true)
		} else {
			kf := func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) { f() }
			err = keybind.KeyReleaseFun(kf).Connect(X, X.RootWin(), btn, true)
		}
	}

	bind("Mod4-v", `notify-send 'hello world!'`, true)
	bind("Mod4-b", `notify-send 'hello world 2!'`, false)

	// Finally, if we want this client to stop responding to key events, we
	// can attach another handler that, when run, detaches all previous
	// handlers.
	// This time, we'll show an example of a KeyRelease binding.
	err = keybind.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			// Use keybind.Detach to detach the root window
			// from all KeyPress *and* KeyRelease handlers.
			keybind.Detach(X, X.RootWin())

			log.Printf("Detached all Key{Press,Release}Events from the "+
				"root window (%d).", X.RootWin())
		}).Connect(X, X.RootWin(), "Mod4-Shift-q", true)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Program initialized. Start pressing keys!")
	xevent.Main(X)
}

func main() {
	btn := flag.String("button", "", "specify a button")
	cmd := flag.String("command", "", "set command for the button")
	clear := flag.Bool("clear", false, "clear the button")
	clearAll := flag.Bool("clearall", false, "clear all bindings")
	flag.Parse()

	if *btn == "" || clearAll == nil {
		log.Println("Starting as a daemon")
		runDaemon()
	} else {
		// TODO
		pipeFile := defaultFIFOPath
		file, err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)
		_ = file
		_ = err
		fmt.Println(*clearAll, *btn, *cmd, *clear)
	}
}
