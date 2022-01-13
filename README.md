# goxhkd

A simple X11 hotkey daemon using [xgbutil](https://github.com/BurntSushi/xgbutil).

The config file is a shell script. It should contain a shebang and be executable. It's executed by the goxhkd daemon on start after a small delay. The following paths are checked:

- `$XDG_CONFIG_HOME/goxhkd/goxhkdrc` (if `$XDG_CONFIG_HOME` isn't set, it falls back to `$HOME/.config`)
- `$HOME/.goxhkdrc`
- `/etc/goxhkd/goxhkdrc`

## Building

- `make build` will build and place the binaries into the project root.
- `make install` will install them to `$GOPATH/bin`.

## Example config

```sh
#!/bin/sh
goxhkc -clearall

for i in $(seq 9); do
      goxhkc -button Mod4-$i bspc desktop -f "^$i"
      goxhkc -button Mod4-Shift-$i bspc node -d "^$i"
done

goxhkc -button Mod4-Q bspc node -c
goxhkc -button Mod4-Return xterm
goxhkc -button Mod4-W notify-send "text" "subtext"
```

## Controller

```txt
Usage of goxhkc:
  -address string
        specify connection address (socket path, host, ...) (default "/tmp/goxhkd.sock")
  -button string
        specify a button
  -clearall
        clear all bindings
  -multi
        allow for multiple bindings to the same button
  -network string
        specify connection network (unix, tcp, ...) (default "unix")
  -onrelease
        run command on button release
  -repeat
        repeatedly run command while the button is pressed
  -sh
        run command with 'sh -c ...'
  -version
        print version and exit
  -window uint
        specify a window
```

## Daemon

```txt
Usage of goxhkd:
  -address string
        specify connection address (socket path, host, ...) (default "/tmp/goxhkd.sock")
  -network string
        specify connection network (unix, tcp, ...) (default "unix")
  -version
        print version and exit
```
