package none

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/google/shlex"
)

// A binding belongs to a group.
// Binding groups are organized into layouts.
// Groups in a layout are ordered, in case of conflicts, higher-order binding is used.

type layout []group

type group []binding

type binding struct {
	command  string
	keycodes []xproto.Keycode
	mods     uint16
	keys     string
}

func (g *group) grabGroup(xu *xgbutil.XUtil) {
	for binding := range *g {
		fmt.Println(binding)
	}
}

func makeCommand(cmd string) (*exec.Cmd, error) {
	tokens, err := shlex.Split(cmd)
	return exec.Command(tokens[0], tokens[1:]...), err
}

func makeBinding(xu *xgbutil.XUtil, str string, cmd string) (b binding, err error) {
	mods, keycodes, err := keybind.ParseString(xu, str)
	if err != nil {
		return
	}
	b = binding{command: cmd, keycodes: keycodes, mods: mods, keys: str}
	return
}

func xorgLoop(xu *xgbutil.XUtil) {
	for {
		ev, err := xu.Conn().WaitForEvent()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ev)
	}
}

func main() {
	xu, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(xu)

	b, err := makeBinding(xu, "mod4-V", "notify-send 'hello5235'")
	g := group{b}
	// b, err = makeBinding(xu, "mod4-V", "notify-send 'hello5235'")
	// g = append(g, b)

	for _, binding := range g {
		callback := func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Println(e)
			cmd, err := makeCommand(binding.command)
			cmd.Start()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		xevent.KeyPressFun(callback).Connect(xu, xu.RootWin())
	}

	var done chan int
	go func(done chan int) {
		xevent.Main(xu)
	}(done)
	<-done
}
