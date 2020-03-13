package shared

const DefaultSocketAddr = "/tmp/goxhkd.sock"

type Connection struct {
	Network string
	Address string
}

func DefaultSocketConnection() *Connection {
	return &Connection{
		Network: "unix",
		Address: DefaultSocketAddr,
	}
}

type Binding struct {
	Cmd          string
	Btn          string
	RunOnRelease bool
	Repeating    bool
}
