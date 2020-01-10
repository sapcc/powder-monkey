package dynomite

// Dynomite represents a dynomite instance
type Dynomite struct {
	Host string
	Port int16
}

// NewDynomite returns a new instance of Dynomite
func NewDynomite(host string, port int16) *Dynomite {
	return &Dynomite{
		Host: host,
		Port: port,
	}
}
