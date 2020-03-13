package main

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xevent"
)

type timedKeyEvent interface {
	GetTime() xproto.Timestamp
	GetKeycode() xproto.Keycode
}

type KeyPressEvent xevent.KeyPressEvent

func (e KeyPressEvent) GetKeycode() xproto.Keycode {
	return e.KeyPressEvent.Detail
}

func (e KeyPressEvent) GetTime() xproto.Timestamp {
	return e.Time
}

type KeyReleaseEvent xevent.KeyReleaseEvent

func (e KeyReleaseEvent) GetKeycode() xproto.Keycode {
	return e.KeyReleaseEvent.Detail
}

func (e KeyReleaseEvent) GetTime() xproto.Timestamp {
	return e.Time
}
