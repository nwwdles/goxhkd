# goxhkd

A simple X11 hotkey daemon using [xgbutil](https://github.com/BurntSushi/xgbutil).

The config file is a shell script (see `examples/`). It should contain a shebang and be executable. It's executed by the goxhkd daemon on start after a small delay. Next config paths are checked:

- `$XDG_CONFIG_HOME/goxhkd/goxhkdrc` (if `$XDG_CONFIG_HOME` isn't set, it falls back to `$HOME/.config`)
- `$HOME/.goxhkdrc`
- `/etc/goxhkd/goxhkdrc`

## Controller

```txt
Usage of goxhkc:
  -address string
        specify connection address (socket path, host, ...) (default "/tmp/goxhkd.sock")
  -button string
        specify a button
  -clear
        clear the button
  -clearall
        clear all bindings
  -command string
        set command for the button
  -network string
        specify connection network (unix, tcp, ...) (default "unix")
  -onrelease
        run command on button release
  -repeat
        repeatedly run command while the button is pressed
  -sh
        run command with 'sh -c ...'
```

## Daemon

```txt
Usage of goxhkd:
  -address string
        specify connection address (socket path, host, ...) (default "/tmp/goxhkd.sock")
  -network string
        specify connection network (unix, tcp, ...) (default "unix")
```
