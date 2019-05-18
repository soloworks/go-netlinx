package apw

// Transport is used to pass connection data around
type Transport struct {
	Type     string
	Host     string
	Port     int
	Name     string
	PingTest bool
	Username string
	Password string
}

// NewIPTransport returns a new project instance with
// default fields already populated
func NewIPTransport(host string) *Transport {
	return &Transport{
		Type: "TCPIP",
		Host: host,
		Port: 1319,
	}
}
