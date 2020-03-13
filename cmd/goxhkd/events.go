package main

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xevent"
)

type timedKeyEvent interface {
	Timestamp() xproto.Timestamp
	Keycode() xproto.Keycode
}

type KeyPressEvent xevent.KeyPressEvent

func (e KeyPressEvent) Keycode() xproto.Keycode {
	return e.KeyPressEvent.Detail
}

func (e KeyPressEvent) Timestamp() xproto.Timestamp {
	return e.Time
}

type KeyReleaseEvent xevent.KeyReleaseEvent

func (e KeyReleaseEvent) Keycode() xproto.Keycode {
	return e.KeyReleaseEvent.Detail
}

func (e KeyReleaseEvent) Timestamp() xproto.Timestamp {
	return e.Time
}
