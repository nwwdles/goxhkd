package main

type Keymap interface {
	ClearKey(key) error
	BindKey(key, cmd string, onPress, noRepeat bool) error
}

func New() Keymap {

}
