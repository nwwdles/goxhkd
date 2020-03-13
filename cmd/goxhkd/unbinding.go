package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

func unbindAll(x *xgbutil.XUtil) error {
	keybind.Detach(x, x.RootWin())
	return nil
}

func unbind(xu *xgbutil.XUtil, btn string, onRelease bool) error {
	mod, codes, err := keybind.ParseString(xu, btn)
	if err != nil {
		err = fmt.Errorf(": %w", err)
		log.Printf("%v", err)

		return err
	}

	code := codes[0]

	switch {
	case onRelease:
		detach(xu, xevent.KeyRelease, xu.RootWin(), mod, code)
	case !onRelease:
		detach(xu, xevent.KeyPress, xu.RootWin(), mod, code)
	}

	return nil
}

// taken from xgbutil
func detach(xu *xgbutil.XUtil, evtype int, win xproto.Window, mods uint16,
	keycode xproto.Keycode) {
	xu.KeybindsLck.RLock()
	defer xu.KeybindsLck.RUnlock()

	for key := range xu.Keybinds {
		if key.Win != win || key.Code != keycode ||
			evtype != key.Evtype || mods != key.Mod {
			continue
		}

		ungrab(xu, key)
	}
}

func ungrab(xu *xgbutil.XUtil, key xgbutil.KeyKey) {
	xu.Keygrabs[key] -= len(xu.Keybinds[key])
	delete(xu.Keybinds, key)

	if xu.Keygrabs[key] == 0 {
		keybind.Ungrab(xu, key.Win, key.Mod, key.Code)
	}
}
