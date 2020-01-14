package dynomite

// Dynomite represents a dynomite instance
type Dynomite struct {
	Host    string
	Port    int16
	Backend *Redis
}

// NewDynomite returns a new instance of Dynomite
func NewDynomite(host string, port int16) *Dynomite {
	return &Dynomite{
		Host: host,
		Port: port,
	}
}

// NewDynomiteRedis creates a new Redis struct with initialized ConnectionPool
func NewDynomiteRedis(host string, port, backendPort int16, password string) *Dynomite {
	redis := NewRedis(host, backendPort, password)

	return &Dynomite{
		Host:    host,
		Port:    port,
		Backend: redis,
	}
}
