package shared

const DefaultSocketAddr = "/tmp/goxhkd.sock"

type Connection struct {
	Network string
	Address string
}

type Binding struct {
	Cmd        string
	Btn        string
	RunOnPress bool
	Repeating  bool
}
